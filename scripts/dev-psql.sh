#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/dev-env.sh"

usage() {
	cat <<'EOF'
Usage:
  ./scripts/dev-psql.sh
  ./scripts/dev-psql.sh "SELECT COUNT(*) FROM note;"
  ./scripts/dev-psql.sh -c "SELECT COUNT(*) FROM note;"
  ./scripts/dev-psql.sh -f /tmp/query.sql
  ./scripts/dev-psql.sh -- --tuples-only -c "SELECT current_database();"

Connects to the local development database from the required .dev file.
Password is passed through PGPASSWORD and is never printed.
EOF
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
	usage
	exit 0
fi

if ! command -v psql >/dev/null 2>&1; then
	echo "psql is required but was not found in PATH" >&2
	exit 1
fi

psql_args=(
	-h "$DB_HOST"
	-p "$DB_PORT"
	-U "$DB_USER"
	-d "$DB_NAME"
	-v ON_ERROR_STOP=1
	-P pager=off
)

export PGPASSWORD="$DB_PASSWORD"

if [[ "${1:-}" == "--target" ]]; then
	echo "${DB_USER}@${DB_HOST}:${DB_PORT}/${DB_NAME} (${APP_DEV_ENV_FILE})"
	exit 0
fi

if [[ $# -eq 0 ]]; then
	echo "Connecting to ${DB_USER}@${DB_HOST}:${DB_PORT}/${DB_NAME}" >&2
	echo "Config: ${APP_DEV_ENV_FILE}" >&2
	exec psql "${psql_args[@]}"
fi

if [[ "$1" == "--" ]]; then
	shift
	exec psql "${psql_args[@]}" "$@"
fi

if [[ "$1" == "-f" || "$1" == "--file" ]]; then
	if [[ -z "${2:-}" ]]; then
		echo "$1 requires a SQL file path" >&2
		exit 1
	fi
	exec psql "${psql_args[@]}" -f "$2"
fi

if [[ "$1" == "-c" || "$1" == "--command" ]]; then
	mode="$1"
	shift
	if [[ $# -eq 0 ]]; then
		echo "$mode requires SQL text" >&2
		exit 1
	fi
	exec psql "${psql_args[@]}" -c "$*"
fi

exec psql "${psql_args[@]}" -c "$*"
