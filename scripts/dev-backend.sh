#!/bin/bash
# Dev backend server script.
# Stops only the configured backend port so another checkout can keep running.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BACKEND_DIR="$PROJECT_ROOT/backend"

if [ -f "$BACKEND_DIR/.env" ]; then
	set -a
	source "$BACKEND_DIR/.env"
	set +a
fi
source "$SCRIPT_DIR/dev-env.sh"

if command -v lsof >/dev/null 2>&1; then
	PIDS="$(lsof -tiTCP:"$BACKEND_PORT" -sTCP:LISTEN 2>/dev/null || true)"
	if [ -n "$PIDS" ]; then
		if [ -t 1 ]; then
			echo "Stopping existing backend listener on port $BACKEND_PORT..."
		fi
		kill $PIDS 2>/dev/null || true
		sleep 0.5
	fi
fi

LOG_DIR="$PROJECT_ROOT/logs"
mkdir -p "$LOG_DIR"
LOG_FILE="$LOG_DIR/backend-${BACKEND_PORT}.log"
> "$LOG_FILE"

export LOG_LEVEL="${1:-TRACE}"
export GOCACHE="${GOCACHE:-/tmp/go-build-cache}"

write_banner() {
	local line="$1"
	if [ -t 1 ]; then
		echo "$line"
	fi
	echo "$line" >> "$LOG_FILE"
}

write_banner "Backend: http://${DEV_BACKEND_HOST}:${BACKEND_PORT}"
write_banner "Database: ${DB_USER}@${DB_HOST}:${DB_PORT}/${DB_NAME}"
write_banner "Config: ${APP_DEV_ENV_FILE}"
write_banner "Log: ${LOG_FILE}"

cd "$BACKEND_DIR" && exec go run ./pservice/entry/ >> "$LOG_FILE" 2>&1
