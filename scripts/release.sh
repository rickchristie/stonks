#!/bin/bash
# Generic production deploy skeleton.
# Customize SSH host, service name, production paths, and reverse proxy before real use.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BUILD_DIR="$PROJECT_ROOT/release/.build"
SSH_HOST="${APP_PROD_SSH_HOST:-app-template-prod-1}"
APP_ROOT="${APP_PROD_ROOT:-/opt/app-template}"
SERVICE_NAME="${APP_PROD_SERVICE_NAME:-app-template-backend}"
VERSION="$(tr -d '\n' < "$PROJECT_ROOT/VERSION")"

usage() {
	cat <<'EOF'
Usage:
  ./scripts/release.sh [--skip-frontend] [--skip-backend]

Environment:
  APP_PROD_SSH_HOST       SSH host, default app-template-prod-1
  APP_PROD_ROOT           Production root, default /opt/app-template
  APP_PROD_SERVICE_NAME   Systemd service, default app-template-backend
EOF
}

SKIP_FRONTEND=false
SKIP_BACKEND=false

while [[ $# -gt 0 ]]; do
	case "$1" in
		--skip-frontend)
			SKIP_FRONTEND=true
			shift
			;;
		--skip-backend)
			SKIP_BACKEND=true
			shift
			;;
		--help|-h)
			usage
			exit 0
			;;
		*)
			echo "Unknown option: $1" >&2
			usage >&2
			exit 1
			;;
	esac
done

ssh -q "$SSH_HOST" exit

rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

if [ "$SKIP_BACKEND" = false ]; then
	cd "$PROJECT_ROOT/backend"
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$VERSION" -o "$BUILD_DIR/app-template-backend" ./pservice/entry/
fi

if [ "$SKIP_FRONTEND" = false ]; then
	cd "$PROJECT_ROOT/web"
	VITE_API_BASE="" VITE_APP_VERSION="$VERSION" npm run build
	mkdir -p "$BUILD_DIR/web"
	cp -r build/* "$BUILD_DIR/web/"
fi

if [ "$SKIP_BACKEND" = false ]; then
	ssh "$SSH_HOST" "sudo systemctl stop $SERVICE_NAME 2>/dev/null || true"
	scp "$BUILD_DIR/app-template-backend" "$SSH_HOST:/tmp/app-template-backend"
	ssh "$SSH_HOST" "sudo mkdir -p $APP_ROOT/backend && sudo mv /tmp/app-template-backend $APP_ROOT/backend/app-template-backend && sudo chmod +x $APP_ROOT/backend/app-template-backend"
fi

if [ "$SKIP_FRONTEND" = false ]; then
	cd "$BUILD_DIR"
	tar -czf web.tar.gz -C web .
	scp web.tar.gz "$SSH_HOST:/tmp/app-template-web.tar.gz"
	ssh "$SSH_HOST" "sudo mkdir -p $APP_ROOT/web && sudo rm -rf $APP_ROOT/web/* && sudo tar -xzf /tmp/app-template-web.tar.gz -C $APP_ROOT/web && rm /tmp/app-template-web.tar.gz"
fi

if [ "$SKIP_BACKEND" = false ]; then
	ssh "$SSH_HOST" "sudo systemctl start $SERVICE_NAME"
fi

rm -rf "$BUILD_DIR"
echo "Deployed version $VERSION to $SSH_HOST"
