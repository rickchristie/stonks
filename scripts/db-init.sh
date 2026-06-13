#!/bin/bash
# Local database setup for development/testing.
# Run with: ./scripts/db-init.sh
# Can be re-run to reset the configured local database.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
SCHEMA_DIR="$REPO_ROOT/backend/accessor/db/pg/setup/schema"
source "$REPO_ROOT/scripts/dev-env.sh"

case "$DB_HOST" in
	localhost|127.0.0.1|::1|/*)
		if [[ "$DB_HOST" = /* ]]; then
			ADMIN_PSQL_TARGET=(-h "$DB_HOST" -p "$DB_PORT" -d postgres)
		else
			ADMIN_PSQL_TARGET=(-p "$DB_PORT" -d postgres)
		fi
		;;
	*)
		echo "Error: db-init.sh can only reset local PostgreSQL targets because it uses sudo -u postgres." >&2
		echo "Configured target is $DB_USER@$DB_HOST:$DB_PORT/$DB_NAME." >&2
		exit 1
		;;
esac

admin_psql() {
	sudo -u postgres psql "${ADMIN_PSQL_TARGET[@]}" "$@"
}

echo "Terminating connections to $DB_NAME..."
admin_psql -v ON_ERROR_STOP=1 -v db_name="$DB_NAME" <<'SQL'
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = :'db_name'
  AND pid <> pg_backend_pid();
SQL

echo "Dropping database $DB_NAME if exists..."
admin_psql -v ON_ERROR_STOP=1 -v db_name="$DB_NAME" <<'SQL'
SELECT format('DROP DATABASE IF EXISTS %I', :'db_name') \gexec
SQL

echo "Ensuring user $DB_USER exists..."
admin_psql -v ON_ERROR_STOP=1 -v db_user="$DB_USER" -v db_password="$DB_PASSWORD" <<'SQL'
SELECT CASE
	WHEN EXISTS (SELECT 1 FROM pg_roles WHERE rolname = :'db_user')
		THEN format('ALTER USER %I WITH PASSWORD %L', :'db_user', :'db_password')
	ELSE format('CREATE USER %I WITH PASSWORD %L', :'db_user', :'db_password')
END \gexec
SQL

echo "Creating database $DB_NAME..."
admin_psql -v ON_ERROR_STOP=1 -v db_name="$DB_NAME" -v db_user="$DB_USER" <<'SQL'
SELECT format('CREATE DATABASE %I OWNER %I', :'db_name', :'db_user') \gexec
SQL

echo "Running schema init scripts..."
for sql_file in "$SCHEMA_DIR"/*.init.sql; do
	if [ -f "$sql_file" ]; then
		echo "Running $(basename "$sql_file")..."
		PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$sql_file"
	fi
done

echo ""
echo "Local database setup complete!"
echo "Connection string: postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
