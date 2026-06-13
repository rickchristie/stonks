#!/bin/bash
# Clean Vite/SvelteKit caches, then start the dev web server through the normal path.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
exec "$SCRIPT_DIR/dev-web.sh" --clean "$@"
