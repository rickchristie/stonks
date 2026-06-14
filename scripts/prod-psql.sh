#!/bin/bash
# Run psql against production without exposing credentials locally.

set -euo pipefail

SSH_HOST="${APP_PROD_SSH_HOST:-stonks-prod-1}"
APP_ROOT="${APP_PROD_ROOT:-/opt/stonks}"

usage() {
	cat <<'EOF'
Usage:
  ./scripts/prod-psql.sh [psql-args...]
  ./scripts/prod-psql.sh --file query.sql [psql-args...]
  cat query.sql | ./scripts/prod-psql.sh [psql-args...]
EOF
}

if [[ "${1:-}" == "--help" ]]; then
	usage
	exit 0
fi

ssh -q "$SSH_HOST" exit

LOCAL_SQL_FILE=""
PSQL_ARGS=()

while [[ $# -gt 0 ]]; do
	case "$1" in
		-f|--file)
			LOCAL_SQL_FILE="${2:-}"
			shift 2
			;;
		--file=*)
			LOCAL_SQL_FILE="${1#*=}"
			shift
			;;
		*)
			PSQL_ARGS+=("$1")
			shift
			;;
	esac
done

if [[ -n "$LOCAL_SQL_FILE" && ! -r "$LOCAL_SQL_FILE" ]]; then
	echo "Local SQL file is not readable: $LOCAL_SQL_FILE" >&2
	exit 1
fi

REMOTE_PSQL_ARGS=("${PSQL_ARGS[@]}")
if [[ -n "$LOCAL_SQL_FILE" || -p /dev/stdin ]]; then
	REMOTE_PSQL_ARGS=(-f - "${REMOTE_PSQL_ARGS[@]}")
fi

QUOTED_REMOTE_ARGS=""
for arg in "${REMOTE_PSQL_ARGS[@]}"; do
	printf -v quoted '%q' "$arg"
	QUOTED_REMOTE_ARGS+=" $quoted"
done

REMOTE_CMD="set -euo pipefail; source <(sudo cat $APP_ROOT/.db_credentials); export PGPASSWORD=\"\$DB_PASSWORD\"; exec psql -h localhost -U \"\$DB_USER\" -d \"\$DB_NAME\"${QUOTED_REMOTE_ARGS}"
printf -v REMOTE_CMD_QUOTED '%q' "$REMOTE_CMD"
SSH_REMOTE_COMMAND="bash -lc $REMOTE_CMD_QUOTED"

if [[ -n "$LOCAL_SQL_FILE" ]]; then
	ssh "$SSH_HOST" "$SSH_REMOTE_COMMAND" < "$LOCAL_SQL_FILE"
else
	ssh "$SSH_HOST" "$SSH_REMOTE_COMMAND"
fi
