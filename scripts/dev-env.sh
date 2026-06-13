#!/bin/bash
# Shared local dev configuration loader.
#
# Source this file before reading DB, port, or Playwright settings.
# `.dev` is required so every checkout owns its ports and database explicitly.

APP_DEV_ENV_SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP_DEV_ENV_PROJECT_ROOT="$(dirname "$APP_DEV_ENV_SCRIPT_DIR")"

if [ "${APP_DEV_ENV_LOADED_ROOT:-}" = "$APP_DEV_ENV_PROJECT_ROOT" ]; then
	return 0 2>/dev/null || exit 0
fi

export PROJECT_ROOT="$APP_DEV_ENV_PROJECT_ROOT"

APP_DEV_ENV_FILE="$PROJECT_ROOT/.dev"
if [ ! -f "$APP_DEV_ENV_FILE" ]; then
	echo "Missing required .dev file at $APP_DEV_ENV_FILE" >&2
	echo "Copy .dev.example to .dev, then set unique ports and DB_NAME for this checkout." >&2
	return 1 2>/dev/null || exit 1
fi

set -a
source "$APP_DEV_ENV_FILE"
set +a

required_dev_env() {
	local key="$1"
	if [ -z "${!key+x}" ] || [ -z "${!key}" ]; then
		echo "Invalid .dev: $key is required and cannot be empty" >&2
		return 1
	fi
}

required_port() {
	local key="$1"
	required_dev_env "$key" || return 1
	if ! [[ "${!key}" =~ ^[0-9]+$ ]]; then
		echo "Invalid .dev: $key must be a numeric port, got ${!key}" >&2
		return 1
	fi
}

required_dev_env DEV_BACKEND_HOST || return 1 2>/dev/null || exit 1
required_port DEV_BACKEND_PORT || return 1 2>/dev/null || exit 1
required_dev_env DEV_WEB_HOST || return 1 2>/dev/null || exit 1
required_port DEV_WEB_PORT || return 1 2>/dev/null || exit 1
required_dev_env DB_HOST || return 1 2>/dev/null || exit 1
required_port DB_PORT || return 1 2>/dev/null || exit 1
required_dev_env DB_USER || return 1 2>/dev/null || exit 1
required_dev_env DB_PASSWORD || return 1 2>/dev/null || exit 1
required_dev_env DB_NAME || return 1 2>/dev/null || exit 1
required_dev_env JWT_SECRET || return 1 2>/dev/null || exit 1
required_dev_env CORS_ALLOWED_ORIGINS || return 1 2>/dev/null || exit 1
required_dev_env VITE_API_BASE || return 1 2>/dev/null || exit 1
required_dev_env PLAYWRIGHT_BASE_URL || return 1 2>/dev/null || exit 1
required_dev_env PLAYWRIGHT_API_BASE || return 1 2>/dev/null || exit 1

export DEV_BACKEND_HOST
export DEV_BACKEND_PORT
export DEV_WEB_HOST
export DEV_WEB_PORT

export BACKEND_HOST="${BACKEND_HOST:-}"
export BACKEND_PORT="${BACKEND_PORT:-$DEV_BACKEND_PORT}"

export DB_HOST
export DB_PORT
export DB_USER
export DB_PASSWORD
export DB_NAME

export JWT_SECRET
export CORS_ALLOWED_ORIGINS
export VITE_API_BASE
export PLAYWRIGHT_BASE_URL
export PLAYWRIGHT_API_BASE

# Vite only exposes VITE_* values to browser code. Derive this from the
# resolved root each run so copied `.dev` files cannot point docs at a sibling checkout.
export APP_WORKSPACE_DIR="$PROJECT_ROOT"
export VITE_APP_WORKSPACE_DIR="$APP_WORKSPACE_DIR"

export APP_DEV_ENV_FILE
export APP_DEV_ENV_LOADED=1
export APP_DEV_ENV_LOADED_ROOT="$PROJECT_ROOT"
