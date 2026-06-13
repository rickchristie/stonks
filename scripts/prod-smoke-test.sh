#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/prod-smoke.sh"

assert_matches() {
	local name="$1"
	local line="$2"
	local matches
	matches="$(printf '%s\n' "$line" | find_failure_logs)"
	if [[ "$matches" != *"$line"* ]]; then
		echo "$name did not match failure logs" >&2
		printf 'line: %s\nmatches: %s\n' "$line" "$matches" >&2
		exit 1
	fi
}

assert_clean() {
	local name="$1"
	local line="$2"
	local matches
	matches="$(printf '%s\n' "$line" | find_failure_logs)"
	if [[ -n "$matches" ]]; then
		echo "$name unexpectedly matched failure logs" >&2
		printf 'line: %s\nmatches: %s\n' "$line" "$matches" >&2
		exit 1
	fi
}

assert_matches "json error level" '{"level":"error","msg":"request failed"}'
assert_matches "json fatal level" '{"level": "fatal", "msg":"startup failed"}'
assert_matches "template error level" 'level=error logger=app-template msg="request failed"'
assert_matches "template fatal level with journal prefix" 'Jun 14 host app-template[1]: level=fatal logger=app-template msg="startup failed"'
assert_matches "legacy bracket error" '[ERROR] request failed'
assert_clean "template info level" 'level=info logger=app-template msg="request ok"'
assert_clean "embedded level text" 'somelevel=error logger=app-template msg="not the structured level field"'

echo "prod smoke log matching tests passed"
