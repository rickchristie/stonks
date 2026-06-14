#!/bin/bash
# Local development health report.
#
# The script gathers bounded evidence for configured backend, web, and database
# reachability. Raw bodies and stderr are kept under /tmp for follow-up debugging.

set -u

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

source "$SCRIPT_DIR/dev-env.sh"

TMP_ROOT="$(mktemp -d "${TMPDIR:-/tmp}/stonks-dev-health.XXXXXX")"

BACKEND_BASE="http://${DEV_BACKEND_HOST}:${BACKEND_PORT}"
WEB_BASE="http://${DEV_WEB_HOST}:${DEV_WEB_PORT}"
BACKEND_LOG="$PROJECT_ROOT/logs/backend-${BACKEND_PORT}.log"
WEB_LOG="$PROJECT_ROOT/logs/web-${DEV_WEB_PORT}.log"

ISSUES=0
WARNINGS=0

section() {
	printf "\n== %s ==\n" "$1"
}

add_issue() {
	ISSUES=$((ISSUES + 1))
	printf "ISSUE: %s\n" "$1"
}

add_warning() {
	WARNINGS=$((WARNINGS + 1))
	printf "WARN: %s\n" "$1"
}

slugify() {
	printf "%s" "$1" | tr '[:upper:]' '[:lower:]' | tr -cs 'a-z0-9' '_'
}

print_file_head() {
	local file="$1"
	local lines="$2"

	if [ ! -s "$file" ]; then
		printf "(empty)\n"
		return
	fi

	sed -n "1,${lines}p" "$file"
}

http_check() {
	local name="$1"
	local expected_code="$2"
	local url="$3"
	local slug
	slug="$(slugify "$name")"

	local body="$TMP_ROOT/${slug}.body"
	local headers="$TMP_ROOT/${slug}.headers"
	local meta="$TMP_ROOT/${slug}.meta"
	local err="$TMP_ROOT/${slug}.err"

	section "$name"
	printf "URL: %s\n" "$url"

	if ! command -v curl >/dev/null 2>&1; then
		add_issue "curl is not available"
		return
	fi

	curl -sS -m 5 -D "$headers" -o "$body" \
		-w "http_code=%{http_code}\ntime_total=%{time_total}\nremote_ip=%{remote_ip}\nsize_download=%{size_download}\ncontent_type=%{content_type}\n" \
		"$url" > "$meta" 2> "$err"
	local curl_status=$?

	if [ "$curl_status" -ne 0 ]; then
		add_issue "curl failed with exit code $curl_status"
		print_file_head "$err" 40
		printf "Artifacts: %s %s %s\n" "$headers" "$body" "$err"
		return
	fi

	print_file_head "$meta" 20

	local http_code
	http_code="$(awk -F= '/^http_code=/{print $2}' "$meta" 2>/dev/null | tail -1)"
	if [ "$http_code" != "$expected_code" ]; then
		add_issue "expected HTTP $expected_code, got ${http_code:-unknown}"
	fi

	printf "Body preview:\n"
	print_file_head "$body" 20
	printf "Artifacts: %s %s\n" "$headers" "$body"
}

port_check() {
	local label="$1"
	local port="$2"
	local found=1

	section "$label listener"
	printf "Port: %s\n" "$port"

	if command -v lsof >/dev/null 2>&1; then
		local out="$TMP_ROOT/${label}-lsof.out"
		lsof -nP -iTCP:"$port" -sTCP:LISTEN > "$out" 2> "$TMP_ROOT/${label}-lsof.err"
		if [ "$?" -eq 0 ] && [ -s "$out" ]; then
			found=0
			print_file_head "$out" 40
		fi
	else
		add_warning "lsof is not available"
	fi

	if [ "$found" -ne 0 ] && command -v ss >/dev/null 2>&1; then
		local out="$TMP_ROOT/${label}-ss.out"
		ss -ltnp > "$out" 2> "$TMP_ROOT/${label}-ss.err"
		awk -v port=":$port" '$0 ~ port { print }' "$out" | sed -n '1,40p'
		if awk -v port=":$port" '$0 ~ port { found=1 } END { exit found ? 0 : 1 }' "$out"; then
			found=0
		fi
	fi

	if [ "$found" -ne 0 ]; then
		add_issue "no listener found on port $port"
	fi
}

database_check() {
	section "database"
	printf "Target: %s@%s:%s/%s\n" "$DB_USER" "$DB_HOST" "$DB_PORT" "$DB_NAME"
	printf "Config: %s\n" "$APP_DEV_ENV_FILE"

	if ! command -v psql >/dev/null 2>&1; then
		add_warning "psql is not available"
		return
	fi

	local out="$TMP_ROOT/database.out"
	local err="$TMP_ROOT/database.err"
	local sql='SELECT current_database() AS db, current_user AS usr, (SELECT COUNT(*) FROM note) AS notes;'

	PGPASSWORD="$DB_PASSWORD" PGCONNECT_TIMEOUT=3 psql \
		-h "$DB_HOST" \
		-p "$DB_PORT" \
		-U "$DB_USER" \
		-d "$DB_NAME" \
		-v ON_ERROR_STOP=1 \
		-P pager=off \
		-c "$sql" > "$out" 2> "$err"
	local status=$?

	if [ "$status" -ne 0 ]; then
		add_issue "database query failed with exit code $status"
		print_file_head "$err" 40
		printf "Artifacts: %s %s\n" "$out" "$err"
		return
	fi

	print_file_head "$out" 40
	printf "Artifacts: %s\n" "$out"
}

process_check() {
	section "processes"

	local out="$TMP_ROOT/processes.out"
	if command -v pgrep >/dev/null 2>&1; then
		pgrep -af 'dev-backend.sh|dev-web.sh|vite|pservice/entry|npm run dev|go run' > "$out" 2> "$TMP_ROOT/processes.err"
		if [ -s "$out" ]; then
			print_file_head "$out" 80
			return
		fi
	fi

	add_warning "no matching dev processes found"
}

tooling_check() {
	section "tooling"

	for tool in curl lsof ss psql node npm go; do
		if command -v "$tool" >/dev/null 2>&1; then
			printf "%s: %s\n" "$tool" "$(command -v "$tool")"
		else
			printf "%s: missing\n" "$tool"
		fi
	done
}

log_tail() {
	local label="$1"
	local file="$2"

	section "$label log"
	printf "Path: %s\n" "$file"

	if [ ! -f "$file" ]; then
		add_warning "$label log does not exist"
		return
	fi

	printf "Size: %s bytes\n" "$(wc -c < "$file")"
	printf "Recent errors:\n"
	tail -200 "$file" | grep -Ei 'panic|fatal|error|failed|listen|vite|ready|local' | tail -40 || true
	printf "\nTail:\n"
	tail -60 "$file"
}

section "config"
printf "Project: %s\n" "$PROJECT_ROOT"
printf "Artifacts: %s\n" "$TMP_ROOT"
printf "Backend: %s\n" "$BACKEND_BASE"
printf "Web: %s\n" "$WEB_BASE"
printf "VITE_API_BASE: %s\n" "$VITE_API_BASE"
printf "PLAYWRIGHT_BASE_URL: %s\n" "$PLAYWRIGHT_BASE_URL"
printf "PLAYWRIGHT_API_BASE: %s\n" "$PLAYWRIGHT_API_BASE"
printf "Config: %s\n" "$APP_DEV_ENV_FILE"

tooling_check
process_check
port_check "backend" "$BACKEND_PORT"
port_check "web" "$DEV_WEB_PORT"
http_check "backend health" "200" "$BACKEND_BASE/health-check"
http_check "backend version" "200" "$BACKEND_BASE/version"
http_check "web prd" "200" "$WEB_BASE/prd"
http_check "web app" "200" "$WEB_BASE/app"
database_check
log_tail "backend" "$BACKEND_LOG"
log_tail "web" "$WEB_LOG"

section "summary"
printf "Artifacts: %s\n" "$TMP_ROOT"
printf "Issues: %s\n" "$ISSUES"
printf "Warnings: %s\n" "$WARNINGS"

if [ "$ISSUES" -ne 0 ]; then
	printf "Status: FAIL\n"
	exit 1
fi

printf "Status: OK\n"
