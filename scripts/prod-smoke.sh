#!/bin/bash
# Non-mutating production smoke checks.

set -euo pipefail

SSH_HOST="${APP_PROD_SSH_HOST:-stonks-prod-1}"
SERVICE_NAME="${APP_PROD_SERVICE_NAME:-stonks-backend}"
EXPECTED_VERSION=""
LOG_LINES=120
# The backend may log as JSON or as key-value text.
# Smoke should fail on either form so a production run cannot hide recent errors.
FAILURE_LOG_PATTERN='("level"[[:space:]]*:[[:space:]]*"(error|fatal)")|(^|[[:space:]])level=(error|fatal)([[:space:]]|$)|panic|fatal|failed to|\[ERROR\]'

usage() {
	cat <<'EOF'
Usage:
  ./scripts/prod-smoke.sh [--expect-version X.Y.Z] [--log-lines N]

Checks SSH, systemd service state, local production /health-check, /version,
and recent logs. This script does not read or print production credentials.
EOF
}

find_failure_logs() {
	grep -Ei "$FAILURE_LOG_PATTERN" || true
}

main() {
	while [[ $# -gt 0 ]]; do
		case "$1" in
			--expect-version)
				EXPECTED_VERSION="${2:-}"
				shift 2
				;;
			--expect-version=*)
				EXPECTED_VERSION="${1#*=}"
				shift
				;;
			--log-lines)
				LOG_LINES="${2:-}"
				shift 2
				;;
			--log-lines=*)
				LOG_LINES="${1#*=}"
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

	local service_state
	service_state="$(ssh "$SSH_HOST" "sudo systemctl is-active $SERVICE_NAME")"
	if [[ "$service_state" != "active" ]]; then
		echo "$SERVICE_NAME is $service_state" >&2
		exit 1
	fi

	local health_resp
	health_resp="$(ssh "$SSH_HOST" 'curl -sS -i http://localhost:8080/health-check')"
	if [[ "$health_resp" != HTTP/*" 200 "* || "$health_resp" != *"Gesundheit!"* ]]; then
		echo "health check failed" >&2
		printf '%s\n' "$health_resp" >&2
		exit 1
	fi

	local version_resp
	version_resp="$(ssh "$SSH_HOST" 'curl -sS http://localhost:8080/version')"
	if [[ -n "$EXPECTED_VERSION" && "$version_resp" != *"\"version\":\"$EXPECTED_VERSION\""* ]]; then
		echo "version mismatch; expected $EXPECTED_VERSION, got $version_resp" >&2
		exit 1
	fi

	local journal_log
	journal_log="$(ssh "$SSH_HOST" "sudo journalctl -u $SERVICE_NAME -n $LOG_LINES --no-pager")"
	local matches
	matches="$(printf '%s\n' "$journal_log" | find_failure_logs)"
	if [[ -n "$matches" ]]; then
		echo "recent backend logs contain failure signatures" >&2
		printf '%s\n' "$matches" >&2
		exit 1
	fi

	echo "production smoke checks passed"
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
	main "$@"
fi
