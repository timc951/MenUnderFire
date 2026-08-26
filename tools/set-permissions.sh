#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="${SCRIPT_DIR}/../backend_go/.env"
CONFIG_FILE="${1:-${SCRIPT_DIR}/test-permissions.json}"

# --- Load .env ---
if [ ! -f "$ENV_FILE" ]; then
  echo "ERROR: .env file not found at $ENV_FILE"
  exit 1
fi

export $(grep -v '^#' "$ENV_FILE" | grep -v '^\s*$' | xargs)

DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"
DB_NAME="${DB_NAME:-menunderfire}"

export PGPASSWORD="$DB_PASSWORD"

psql_cmd() {
  psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -A "$@"
}

# --- Load JSON config ---
if [ ! -f "$CONFIG_FILE" ]; then
  echo "ERROR: Config file not found at $CONFIG_FILE"
  exit 1
fi

echo "Reading config from: $CONFIG_FILE"

EMAIL=$(jq -r '.email' "$CONFIG_FILE")
SITE_ADMIN=$(jq -r '.site_admin' "$CONFIG_FILE")

if [ "$EMAIL" = "null" ] || [ "$EMAIL" = "YOUR_EMAIL_HERE" ]; then
  echo "ERROR: Set your email in $CONFIG_FILE"
  exit 1
fi

# --- Resolve user ID ---
USER_ID=$(psql_cmd -c "SELECT id FROM users WHERE email = '$EMAIL'" 2>/dev/null || true)
if [ -z "$USER_ID" ]; then
  echo "ERROR: No user found with email '$EMAIL'"
  exit 1
fi
echo "User: $EMAIL ($USER_ID)"
echo "---"

# --- Site admin ---
if [ "$SITE_ADMIN" = "true" ]; then
  psql_cmd -c "UPDATE users SET is_site_admin = true WHERE id = '$USER_ID'" > /dev/null
  echo "Site admin: ENABLED"
elif [ "$SITE_ADMIN" = "false" ]; then
  psql_cmd -c "UPDATE users SET is_site_admin = false WHERE id = '$USER_ID'" > /dev/null
  echo "Site admin: DISABLED"
fi

# --- Org admin ---
ORG_COUNT=$(jq '.org_admin | length' "$CONFIG_FILE")
for (( i=0; i<ORG_COUNT; i++ )); do
  ACTION=$(jq -r ".org_admin[$i].action" "$CONFIG_FILE")
  ORG_NAME=$(jq -r ".org_admin[$i].org_name" "$CONFIG_FILE")

  ORG_ID=$(psql_cmd -c "SELECT id FROM organizations WHERE name = '$ORG_NAME'" 2>/dev/null || true)
  if [ -z "$ORG_ID" ]; then
    echo "Org admin: SKIPPED - org '$ORG_NAME' not found"
    continue
  fi

  if [ "$ACTION" = "add" ]; then
    psql_cmd -c "
      INSERT INTO organization_admins (user_id, organization_id, invited_by_id)
      VALUES ('$USER_ID', '$ORG_ID', '$USER_ID')
      ON CONFLICT (user_id, organization_id) DO NOTHING
    " > /dev/null
    echo "Org admin: ADDED to '$ORG_NAME'"
  elif [ "$ACTION" = "remove" ]; then
    psql_cmd -c "
      DELETE FROM organization_admins
      WHERE user_id = '$USER_ID' AND organization_id = '$ORG_ID'
    " > /dev/null
    echo "Org admin: REMOVED from '$ORG_NAME'"
  fi
done

# --- Group memberships ---
GROUP_COUNT=$(jq '.group_memberships | length' "$CONFIG_FILE")
for (( i=0; i<GROUP_COUNT; i++ )); do
  ACTION=$(jq -r ".group_memberships[$i].action" "$CONFIG_FILE")
  GROUP_NAME=$(jq -r ".group_memberships[$i].group_name" "$CONFIG_FILE")
  ROLE=$(jq -r ".group_memberships[$i].role // empty" "$CONFIG_FILE")

  GROUP_ID=$(psql_cmd -c "SELECT id FROM groups WHERE name = '$GROUP_NAME'" 2>/dev/null || true)
  if [ -z "$GROUP_ID" ]; then
    echo "Group: SKIPPED - group '$GROUP_NAME' not found"
    continue
  fi

  if [ "$ACTION" = "set" ]; then
    if [ -z "$ROLE" ]; then
      echo "Group: SKIPPED - role is required for 'set' action"
      continue
    fi
    psql_cmd -c "
      INSERT INTO group_memberships (user_id, group_id, role)
      VALUES ('$USER_ID', '$GROUP_ID', '$ROLE')
      ON CONFLICT (user_id, group_id) DO UPDATE SET role = '$ROLE'
    " > /dev/null
    echo "Group: SET $ROLE in '$GROUP_NAME'"
  elif [ "$ACTION" = "remove" ]; then
    psql_cmd -c "
      DELETE FROM group_memberships
      WHERE user_id = '$USER_ID' AND group_id = '$GROUP_ID'
    " > /dev/null
    echo "Group: REMOVED from '$GROUP_NAME'"
  fi
done

# --- Print current status ---
echo ""
echo "=== Current Status ==="
psql_cmd -c "
  SELECT
    u.is_site_admin AS site_admin,
    (SELECT string_agg(o.name, ', ')
     FROM organization_admins oa
     JOIN organizations o ON oa.organization_id = o.id
     WHERE oa.user_id = u.id) AS org_admin_of,
    (SELECT string_agg(gm.role || ' in ' || g.name, ', ')
     FROM group_memberships gm
     JOIN groups g ON gm.group_id = g.id
     WHERE gm.user_id = u.id) AS group_roles
  FROM users u
  WHERE u.id = '$USER_ID'
" | while IFS='|' read -r sa orgs groups; do
  echo "  Site admin:  $sa"
  echo "  Org admin:   ${orgs:-none}"
  echo "  Groups:      ${groups:-none}"
done
