#!/usr/bin/env bash
# DEV ONLY — Sets up the Zitadel OIDC application for local development.
# Idempotent: safe to run multiple times.
# Requires: Docker stack running (make dev-detached)
# Usage: ./infra/scripts/setup-zitadel.sh
set -euo pipefail

ZITADEL_HOST="http://localhost:8080"
ZITADEL_HEADER="Host: zitadel:8080"
ENV_FILE=".env"
PROJECT_NAME="Nutrometra API"
APP_NAME="Nutrometra Web"

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

# --- 3. Find or create OIDC application ---
echo "Looking for app '${APP_NAME}'..."
CLIENT_ID=$(curl -s "${ZITADEL_HOST}/management/v1/projects/${PROJECT_ID}/apps/_search" \
  -H "${ZITADEL_HEADER}" -H "${AUTH}" -H "Content-Type: application/json" \
  -d "{\"queries\":[{\"nameQuery\":{\"name\":\"${APP_NAME}\",\"method\":\"TEXT_QUERY_METHOD_EQUALS\"}}]}" \
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

if [ -n "$CLIENT_ID" ]; then
  echo "App already exists. Client ID: ${CLIENT_ID}"
else
  echo "Creating OIDC application '${APP_NAME}'..."
  CLIENT_ID=$(curl -s -X POST "${ZITADEL_HOST}/management/v1/projects/${PROJECT_ID}/apps/oidc" \
    -H "${ZITADEL_HEADER}" -H "${AUTH}" -H "Content-Type: application/json" \
    -d '{
      "name": "'"${APP_NAME}"'",
      "redirectUris": ["http://localhost:3000/api/auth/callback/zitadel"],
      "postLogoutRedirectUris": ["http://localhost:3000"],
      "responseTypes": ["OIDC_RESPONSE_TYPE_CODE"],
      "grantTypes": ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE", "OIDC_GRANT_TYPE_REFRESH_TOKEN"],
      "appType": "OIDC_APP_TYPE_WEB",
      "authMethodType": "OIDC_AUTH_METHOD_TYPE_NONE",
      "accessTokenType": "OIDC_TOKEN_TYPE_JWT",
      "idTokenRoleAssertion": true,
      "idTokenUserinfoAssertion": true
    }' \
    | python3 -c "import sys,json; print(json.load(sys.stdin)['clientId'])")
  echo "App created. Client ID: ${CLIENT_ID}"
fi

if [ -z "$CLIENT_ID" ]; then
  echo "ERROR: Could not obtain Client ID."
  exit 1
fi

# --- 4. Update .env ---
if [ -f "$ENV_FILE" ]; then
  if grep -q "^ZITADEL_CLIENT_ID=" "$ENV_FILE"; then
    sed -i.bak "s|^ZITADEL_CLIENT_ID=.*|ZITADEL_CLIENT_ID=${CLIENT_ID}|" "$ENV_FILE"
    rm -f "${ENV_FILE}.bak"
  else
    echo "ZITADEL_CLIENT_ID=${CLIENT_ID}" >> "$ENV_FILE"
  fi
  echo "Updated ${ENV_FILE}"
fi

echo ""
echo "=== Zitadel OIDC setup complete ==="
echo "Client ID: ${CLIENT_ID}"
