#!/bin/bash
# Initialize database with small sample data.
# Run after db-init.sh: ./scripts/db-seed.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
source "$REPO_ROOT/scripts/dev-env.sh"

PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<'SQL'
TRUNCATE note RESTART IDENTITY;

INSERT INTO note (title, body, status)
VALUES
	('Hello World', 'This message came from PostgreSQL through the Go backend.', 'Active'),
	('Archived example', 'Archived rows stay available to direct get calls but are hidden from the active list.', 'Archived');
SQL

echo "Test data initialized successfully."
