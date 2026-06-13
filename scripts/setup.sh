#!/bin/bash
# Install template dependencies after copying the repository.
#
# This intentionally does not source .dev; dependency installation should work
# before local ports and database names are finalized. npm still reads each
# project .npmrc, including min-release-age and engine-strict.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT/backend"
go mod download

cd "$PROJECT_ROOT/web"
npm install

cd "$PROJECT_ROOT/web/playtest"
npm install
