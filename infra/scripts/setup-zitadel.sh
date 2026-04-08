#!/usr/bin/env bash
# DEV ONLY — Sets up Zitadel OIDC applications for local development.
# Creates 3 apps: web-professional (3000), web-patient (3001), backoffice (3002).
# Idempotent: safe to run multiple times.
# Requires: Docker stack running (make dev-detached)
# Usage: ./infra/scripts/setup-zitadel.sh
set -euo pipefail

ZITADEL_HOST="http://localhost:8080"
ZITADEL_HEADER="Host: zitadel:8080"
ENV_FILE=".env"
PROJECT_NAME="Nutrometra API"

# --- 1. Get PAT from Zitadel container ---
echo "Extracting PAT from Zitadel container..."
PAT_TMP=$(mktemp)
docker cp infra-zitadel-1:/machinekey/zitadel-admin-sa.pat "$PAT_TMP" 2>/dev/null || {
  echo "ERROR: Could not extract PAT from container."
  rm -f "$PAT_TMP"
  exit 1
}
PAT=$(tr -d '[:space:]' < "$PAT_TMP")
rm -f "$PAT_TMP"

if [ -z "$PAT" ]; then
  echo "ERROR: PAT file is empty."
  exit 1
fi
echo "PAT obtained."

AUTH="Authorization: Bearer ${PAT}"

# --- 2. Find or create project ---
echo "Looking for project '${PROJECT_NAME}'..."
PROJECT_ID=$(curl -s "${ZITADEL_HOST}/management/v1/projects/_search" \
  -H "${ZITADEL_HEADER}" -H "${AUTH}" -H "Content-Type: application/json" \
  -d "{\"queries\":[{\"nameQuery\":{\"name\":\"${PROJECT_NAME}\",\"method\":\"TEXT_QUERY_METHOD_EQUALS\"}}]}" \
  | python3 -c "
import sys, json
data = json.load(sys.stdin)
projects = data.get('result', [])
print(projects[0]['id'] if projects else '')
" 2>/dev/null)

if [ -n "$PROJECT_ID" ]; then
  echo "Project already exists: ${PROJECT_ID}"
else
  echo "Creating project '${PROJECT_NAME}'..."
  PROJECT_ID=$(curl -s -X POST "${ZITADEL_HOST}/management/v1/projects" \
    -H "${ZITADEL_HEADER}" -H "${AUTH}" -H "Content-Type: application/json" \
    -d "{\"name\": \"${PROJECT_NAME}\", \"projectRoleAssertion\": true}" \
    | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")
  echo "Project created: ${PROJECT_ID}"
fi

# --- 3. Helper: find or create OIDC app ---
# Usage: create_oidc_app "App Name" "http://localhost:PORT" "ENV_VAR_NAME"
create_oidc_app() {
  local app_name="$1"
  local base_url="$2"
  local env_var="$3"

  echo ""
  echo "--- Setting up '${app_name}' (${base_url}) ---"

  local client_id
  client_id=$(curl -s "${ZITADEL_HOST}/management/v1/projects/${PROJECT_ID}/apps/_search" \
    -H "${ZITADEL_HEADER}" -H "${AUTH}" -H "Content-Type: application/json" \
    -d "{\"queries\":[{\"nameQuery\":{\"name\":\"${app_name}\",\"method\":\"TEXT_QUERY_METHOD_EQUALS\"}}]}" \
    | python3 -c "
import sys, json
data = json.load(sys.stdin)
apps = data.get('result', [])
for app in apps:
    oidc = app.get('oidcConfig', {})
    cid = oidc.get('clientId', '')
    if cid:
        print(cid)
        break
" 2>/dev/null)

  if [ -n "$client_id" ]; then
    echo "App already exists. Client ID: ${client_id}"
  else
    echo "Creating OIDC application '${app_name}'..."
    client_id=$(curl -s -X POST "${ZITADEL_HOST}/management/v1/projects/${PROJECT_ID}/apps/oidc" \
      -H "${ZITADEL_HEADER}" -H "${AUTH}" -H "Content-Type: application/json" \
      -d '{
        "name": "'"${app_name}"'",
        "redirectUris": ["'"${base_url}"'/api/auth/callback/zitadel"],
        "postLogoutRedirectUris": ["'"${base_url}"'"],
        "responseTypes": ["OIDC_RESPONSE_TYPE_CODE"],
        "grantTypes": ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE", "OIDC_GRANT_TYPE_REFRESH_TOKEN"],
        "appType": "OIDC_APP_TYPE_WEB",
        "authMethodType": "OIDC_AUTH_METHOD_TYPE_NONE",
        "accessTokenType": "OIDC_TOKEN_TYPE_JWT",
        "idTokenRoleAssertion": true,
        "idTokenUserinfoAssertion": true
      }' \
      | python3 -c "import sys,json; print(json.load(sys.stdin)['clientId'])")
    echo "App created. Client ID: ${client_id}"
  fi

  if [ -z "$client_id" ]; then
    echo "ERROR: Could not obtain Client ID for '${app_name}'."
    exit 1
  fi

  # Update .env
  if [ -f "$ENV_FILE" ]; then
    if grep -q "^${env_var}=" "$ENV_FILE"; then
      sed -i.bak "s|^${env_var}=.*|${env_var}=${client_id}|" "$ENV_FILE"
      rm -f "${ENV_FILE}.bak"
    else
      echo "${env_var}=${client_id}" >> "$ENV_FILE"
    fi
    echo "Updated ${ENV_FILE}: ${env_var}=${client_id}"
  fi
}

# --- 4. Create all three OIDC applications ---
create_oidc_app "Nutrometra Web"            "http://localhost:3000" "ZITADEL_CLIENT_ID"
create_oidc_app "Nutrometra Web - Patient"  "http://localhost:3001" "ZITADEL_CLIENT_ID_PATIENT"
create_oidc_app "Nutrometra Web - Backoffice" "http://localhost:3002" "ZITADEL_CLIENT_ID_BACKOFFICE"

echo ""
echo "=== Zitadel OIDC setup complete ==="
echo "3 applications configured for ports 3000, 3001, 3002"
