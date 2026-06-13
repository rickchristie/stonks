#!/bin/bash
# Initialize database with schema and no app rows.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
source "$REPO_ROOT/scripts/dev-env.sh"

PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<'SQL'
TRUNCATE note RESTART IDENTITY;
SQL

echo "Empty data initialized successfully."
