#!/bin/bash
# Guided Playwright runner that keeps command output under logs/.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
PLAYTEST_DIR="$PROJECT_ROOT/web/playtest"
LOG_DIR="$PROJECT_ROOT/logs"
source "$SCRIPT_DIR/dev-env.sh"

EXTRA_ARGS=("$@")
BROWSER_LABEL=""
BROWSER_ARGS=()
WORKER_ARGS=()

shopt -s nullglob globstar

print_header() {
	clear 2>/dev/null || true
	printf 'Playwright Test Runner\n'
	if [ -n "${BROWSER_LABEL:-}" ]; then
		printf 'Mode: %s\n' "$BROWSER_LABEL"
	fi
	printf '\n'
}

collect_directories() {
	local path dir
	DIRECTORIES=()

	for path in "$PLAYTEST_DIR"/*/; do
		dir="$(basename "$path")"
		case "$dir" in
			assets|manual|node_modules|test-results)
				continue
				;;
		esac

		local specs=("$path"**/*.spec.ts)
		if [ ${#specs[@]} -gt 0 ]; then
			DIRECTORIES+=("$dir")
		fi
	done
}

print_numbered_list() {
	local items=("$@")
	local i
	for ((i = 0; i < ${#items[@]}; i++)); do
		printf '%d. %s\n' "$((i + 1))" "${items[$i]}"
	done
}

parse_multi_numbers() {
	local input="$1"
	local max="$2"
	local output_name="$3"
	local -n output_ref="$output_name"
	local token cleaned
	local -A seen=()

	IFS=',' read -r -a tokens <<< "$input"
	output_ref=()

	for token in "${tokens[@]}"; do
		cleaned="${token//[[:space:]]/}"
		if [ -z "$cleaned" ]; then
			continue
		fi
		if [[ ! "$cleaned" =~ ^[0-9]+$ ]]; then
			return 1
		fi
		if (( cleaned < 1 || cleaned > max )); then
			return 1
		fi
		if [ -z "${seen[$cleaned]+x}" ]; then
			seen[$cleaned]=1
			output_ref+=("$cleaned")
		fi
	done

	if [ ${#output_ref[@]} -eq 0 ]; then
		return 1
	fi
}

parse_single_number() {
	local input="${1//[[:space:]]/}"
	local max="$2"
	if [[ ! "$input" =~ ^[0-9]+$ ]]; then
		return 1
	fi
	if (( input < 1 || input > max )); then
		return 1
	fi
	return 0
}

collect_files_for_directory() {
	local dir="$1"
	local path rel
	FILES=()

	for path in "$PLAYTEST_DIR/$dir"/**/*.spec.ts; do
		rel="${path#"$PLAYTEST_DIR/"}"
		FILES+=("$rel")
	done
}

extra_args_have_workers() {
	local arg

	for arg in "${EXTRA_ARGS[@]}"; do
		case "$arg" in
			--workers|--workers=*|-j|-j*)
				return 0
				;;
		esac
	done

	return 1
}

prompt_worker_count() {
	local default_workers="$1"
	local worker_input

	if extra_args_have_workers; then
		WORKER_ARGS=()
		return
	fi

	while true; do
		print_header
		printf 'Choose worker count:\n'
		printf 'Press Enter to use default: %s\n\n' "$default_workers"
		read -r -p 'Enter workers: ' worker_input
		worker_input="${worker_input//[[:space:]]/}"

		if [ -z "$worker_input" ]; then
			WORKER_ARGS=("--workers=$default_workers")
			return
		fi

		if [[ "$worker_input" =~ ^[0-9]+$ ]] && (( worker_input > 0 )); then
			WORKER_ARGS=("--workers=$worker_input")
			return
		fi

		printf '\nInvalid worker count. Press Enter to try again.'
		read -r
	done
}

sanitize_for_log() {
	perl -CSDA -pe 's/\e\[[0-9;?]*[ -\/]*[@-~]//g; s/\e\][^\a]*(?:\a|\e\\)//g; s/\r/\n/g; s/[\x{2800}-\x{28FF}]//g;'
}

run_playwright() {
	local targets=("$@")
	local timestamp
	local log_file
	local default_workers
	local command

	if [ ${#targets[@]} -eq 0 ]; then
		default_workers=8
	else
		default_workers=4
	fi

	prompt_worker_count "$default_workers"

	mkdir -p "$LOG_DIR"
	timestamp="$(date +%Y%m%d-%H%M%S)"
	log_file="$LOG_DIR/playwright-$timestamp.log"

	printf -v command 'cd %q && npx playwright test' "$PLAYTEST_DIR"
	for target in "${targets[@]}"; do
		printf -v command '%s %q' "$command" "$target"
	done
	for arg in "${BROWSER_ARGS[@]}"; do
		printf -v command '%s %q' "$command" "$arg"
	done
	for arg in "${WORKER_ARGS[@]}"; do
		printf -v command '%s %q' "$command" "$arg"
	done
	for arg in "${EXTRA_ARGS[@]}"; do
		printf -v command '%s %q' "$command" "$arg"
	done

	print_header
	printf 'Running:\n  %s' "$command"
	printf '\nLog:\n  %s\n\n' "$log_file"

	printf 'Command: %s\n\n' "$command" > "$log_file"
	if command -v script >/dev/null 2>&1; then
		script -qefc "$command" /dev/null | tee >(sanitize_for_log >> "$log_file")
	else
		FORCE_COLOR=1 bash -lc "$command" 2>&1 | tee >(sanitize_for_log >> "$log_file")
	fi
}

collect_directories

if [ ${#DIRECTORIES[@]} -eq 0 ]; then
	printf 'No Playwright test directories found in %s\n' "$PLAYTEST_DIR" >&2
	exit 1
fi

while true; do
	print_header
	printf 'Choose browser mode:\n'
	printf '1. Headless\n'
	printf '2. Headed\n\n'
	read -r -p 'Enter choice: ' browser_choice
	browser_choice="${browser_choice//[[:space:]]/}"

	case "$browser_choice" in
		1)
			BROWSER_LABEL='Headless'
			BROWSER_ARGS=()
			break
			;;
		2)
			BROWSER_LABEL='Headed'
			BROWSER_ARGS=("--headed")
			break
			;;
		*)
			printf '\nInvalid choice. Press Enter to try again.'
			read -r
			;;
	esac
done

while true; do
	print_header
	printf 'Choose what to run:\n'
	printf '1. Run all playtests\n'
	printf '2. Run directories\n'
	printf '3. Run file\n'
	printf '4. Type in suites\n\n'
	read -r -p 'Enter choice: ' main_choice
	main_choice="${main_choice//[[:space:]]/}"

	case "$main_choice" in
		1)
			run_playwright
			exit 0
			;;
		2)
			while true; do
				print_header
				printf 'Directories:\n\n'
				print_numbered_list "${DIRECTORIES[@]}"
				printf '\n'
				read -r -p 'Enter directory numbers separated by commas: ' dir_input
				if parse_multi_numbers "$dir_input" "${#DIRECTORIES[@]}" selected_dir_indexes; then
					selected_dirs=()
					for index in "${selected_dir_indexes[@]}"; do
						selected_dirs+=("${DIRECTORIES[$((index - 1))]}/")
					done
					run_playwright "${selected_dirs[@]}"
					exit 0
				fi
				printf '\nInvalid selection. Press Enter to try again.'
				read -r
			done
			;;
		3)
			while true; do
				print_header
				printf 'Directories:\n\n'
				print_numbered_list "${DIRECTORIES[@]}"
				printf '\n'
				read -r -p 'Enter one directory number: ' single_dir_input
				if parse_single_number "$single_dir_input" "${#DIRECTORIES[@]}"; then
					selected_directory="${DIRECTORIES[$((single_dir_input - 1))]}"
					break
				fi
				printf '\nInvalid selection. Press Enter to try again.'
				read -r
			done

			collect_files_for_directory "$selected_directory"
			if [ ${#FILES[@]} -eq 0 ]; then
				printf 'No spec files found in %s\n' "$selected_directory" >&2
				exit 1
			fi

			while true; do
				print_header
				printf 'Files in %s:\n\n' "$selected_directory"
				print_numbered_list "${FILES[@]}"
				printf '\n'
				read -r -p 'Enter file numbers separated by commas: ' file_input
				if parse_multi_numbers "$file_input" "${#FILES[@]}" selected_file_indexes; then
					selected_files=()
					for index in "${selected_file_indexes[@]}"; do
						selected_files+=("${FILES[$((index - 1))]}")
					done
					run_playwright "${selected_files[@]}"
					exit 0
				fi
				printf '\nInvalid selection. Press Enter to try again.'
				read -r
			done
			;;
		4)
			while true; do
				print_header
				printf 'Type suites, directories, or file paths separated by spaces.\n'
				printf 'Example: app/ prd/prd-route.spec.ts\n\n'
				read -r -p 'Enter suites: ' suite_input
				read -r -a typed_suites <<< "$suite_input"
				if [ ${#typed_suites[@]} -gt 0 ]; then
					run_playwright "${typed_suites[@]}"
					exit 0
				fi
				printf '\nInvalid selection. Press Enter to try again.'
				read -r
			done
			;;
		*)
			printf '\nInvalid choice. Press Enter to try again.'
			read -r
			;;
	esac
done
