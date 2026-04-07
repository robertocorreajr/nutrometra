#!/usr/bin/env bash
# DEV ONLY — Sets up the Zitadel OIDC application for local development.
# Requires: Docker stack running (make dev-detached)
# Usage: ./infra/scripts/setup-zitadel.sh
set -euo pipefail

ZITADEL_HOST="http://localhost:8080"
ZITADEL_INTERNAL="http://zitadel:8080"
ENV_FILE=".env"

# --- 1. Get PAT from Zitadel container ---
echo "Extracting PAT from Zitadel container..."
PAT_TMP=$(mktemp)
docker cp infra-zitadel-1:/machinekey/zitadel-admin-sa.pat "$PAT_TMP" 2>/dev/null || {
  echo "ERROR: Could not find PAT file in Zitadel container."
  echo "Ensure init.yaml includes Machine + Pat config and /machinekey volume is writable."
  rm -f "$PAT_TMP"
  exit 1
}
PAT=$(tr -d '[:space:]' < "$PAT_TMP")
rm -f "$PAT_TMP"
echo "PAT obtained."

# --- 2. Create project ---
echo "Creating project 'Nutrometra API'..."
# Requests go to localhost:8080 but Zitadel expects Host: zitadel
PROJECT_RESPONSE=$(curl -sf -X POST "${ZITADEL_HOST}/management/v1/projects" \
  -H "Host: zitadel:8080" \
  -H "Authorization: Bearer ${PAT}" \
  -H "Content-Type: application/json" \
  -d '{"name": "Nutrometra API", "projectRoleAssertion": true}') || {
  echo "Failed to create project."
  exit 1
}

PROJECT_ID=$(echo "$PROJECT_RESPONSE" | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")
echo "Project created: ${PROJECT_ID}"

# --- 3. Create OIDC application ---
echo "Creating OIDC application..."
APP_RESPONSE=$(curl -sf -X POST "${ZITADEL_HOST}/management/v1/projects/${PROJECT_ID}/apps/oidc" \
  -H "Host: zitadel:8080" \
  -H "Authorization: Bearer ${PAT}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Nutrometra Web",
    "redirectUris": ["http://localhost:3000/api/auth/callback/zitadel"],
    "postLogoutRedirectUris": ["http://localhost:3000"],
    "responseTypes": ["OIDC_RESPONSE_TYPE_CODE"],
    "grantTypes": ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE", "OIDC_GRANT_TYPE_REFRESH_TOKEN"],
    "appType": "OIDC_APP_TYPE_WEB",
    "authMethodType": "OIDC_AUTH_METHOD_TYPE_NONE",
    "accessTokenType": "OIDC_TOKEN_TYPE_JWT",
    "idTokenRoleAssertion": true,
    "idTokenUserinfoAssertion": true
  }') || {
  echo "Failed to create OIDC app."
  exit 1
}

CLIENT_ID=$(echo "$APP_RESPONSE" | python3 -c "import sys,json; print(json.load(sys.stdin)['clientId'])")
echo "OIDC app created. Client ID: ${CLIENT_ID}"

# --- 4. Update .env ---
if [ -f "$ENV_FILE" ]; then
  if grep -q "^ZITADEL_CLIENT_ID=" "$ENV_FILE"; then
    sed -i.bak "s|^ZITADEL_CLIENT_ID=.*|ZITADEL_CLIENT_ID=${CLIENT_ID}|" "$ENV_FILE"
    rm -f "${ENV_FILE}.bak"
  else
    echo "ZITADEL_CLIENT_ID=${CLIENT_ID}" >> "$ENV_FILE"
  fi
  echo "Updated ${ENV_FILE} with ZITADEL_CLIENT_ID=${CLIENT_ID}"
else
  echo "ZITADEL_CLIENT_ID=${CLIENT_ID}"
  echo "No .env file found. Set ZITADEL_CLIENT_ID=${CLIENT_ID} in your environment."
fi

echo ""
echo "=== Zitadel OIDC setup complete ==="
echo "Client ID: ${CLIENT_ID}"
