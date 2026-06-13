#!/bin/bash
# Dev web server script.
# Stops only the configured web port so another checkout can keep running.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
WEB_DIR="$PROJECT_ROOT/web"

source "$SCRIPT_DIR/dev-env.sh"

if command -v lsof >/dev/null 2>&1; then
	PIDS="$(lsof -tiTCP:"$DEV_WEB_PORT" -sTCP:LISTEN 2>/dev/null || true)"
	if [ -n "$PIDS" ]; then
		if [ -t 1 ]; then
			echo "Stopping existing web listener on port $DEV_WEB_PORT..."
		fi
		kill $PIDS 2>/dev/null || true
		sleep 0.5
	fi
fi

LOG_DIR="$PROJECT_ROOT/logs"
mkdir -p "$LOG_DIR"
LOG_FILE="$LOG_DIR/web-${DEV_WEB_PORT}.log"
> "$LOG_FILE"

write_banner() {
	local line="$1"
	if [ -t 1 ]; then
		echo "$line"
	fi
	echo "$line" >> "$LOG_FILE"
}

write_banner "Web: http://${DEV_WEB_HOST}:${DEV_WEB_PORT}"
write_banner "API: ${VITE_API_BASE}"
write_banner "Workspace: ${VITE_APP_WORKSPACE_DIR}"
write_banner "Config: ${APP_DEV_ENV_FILE}"
write_banner "Log: ${LOG_FILE}"

if [ "${1:-}" = "--clean" ]; then
	cd "$WEB_DIR" && exec npm run dev:clean >> "$LOG_FILE" 2>&1
else
	cd "$WEB_DIR" && exec npm run dev >> "$LOG_FILE" 2>&1
fi
