#!/bin/bash
# Vite dev command used by `cd web && npm run dev`.
# Loading repo-local `.dev` here keeps direct npm usage aligned with scripts/dev-web.sh.

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/dev-env.sh"

if [ "${1:-}" = "--clean" ]; then
	shift
	rm -rf "$PROJECT_ROOT/web/.svelte-kit" "$PROJECT_ROOT/web/node_modules/.vite"
fi

export VITE_APP_VERSION="${VITE_APP_VERSION:-$(tr -d '\n' < "$PROJECT_ROOT/VERSION")}"

exec vite dev --host "$DEV_WEB_HOST" --port "$DEV_WEB_PORT" --strictPort "$@"
