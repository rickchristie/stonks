#!/bin/bash
# Local database initialization wrapper.
# Reads the required .dev file through scripts/dev-env.sh.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/dev-env.sh"

"$SCRIPT_DIR/db-init.sh"
"$SCRIPT_DIR/db-seed.sh"
