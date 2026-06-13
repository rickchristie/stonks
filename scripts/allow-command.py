#!/usr/bin/env python3
"""Shared allow-list checker for agent permission hooks.

The script is tool-neutral:
- Codex can call it from a PermissionRequest hook and receive Codex allow JSON.
- OpenCode or other tools can call `--check -- <command>` and use the exit code.
- Humans can call `--json -- <command>` while editing rules.
"""

from __future__ import annotations

import argparse
import json
from pathlib import Path
import re
import shlex
import sys
from typing import Any


SCRIPT_DIR = Path(__file__).resolve().parent
REPO_ROOT = SCRIPT_DIR.parent
RULES_PATH = SCRIPT_DIR / "allow-command-rules.json"


def main() -> int:
	args = parse_args()
	command = get_command_from_args(args)
	if command is None:
		command = get_command_from_stdin()

	if not command:
		return 1 if args.check else 0

	result = check_command(command, cwd=args.cwd)

	if args.json:
		print(json.dumps(result))
	elif args.codex or command_was_from_stdin(args):
		if result["allowed"]:
			print(json.dumps({
				"hookSpecificOutput": {
					"hookEventName": "PermissionRequest",
					"decision": {"behavior": "allow"},
					"rule": result["rule"],
				}
			}))
	elif not args.check:
		status = "ALLOW" if result["allowed"] else "DENY"
		suffix = f" by {result['rule']}" if result["rule"] else ""
		print(f"{status}{suffix}")

	return 0 if result["allowed"] else 1


def parse_args() -> argparse.Namespace:
	parser = argparse.ArgumentParser(description="Check whether a command is on the repo allow list.")
	parser.add_argument("--check", action="store_true", help="Use exit code only.")
	parser.add_argument("--json", action="store_true", help="Print JSON result.")
	parser.add_argument("--codex", action="store_true", help="Print Codex PermissionRequest allow JSON when allowed.")
	parser.add_argument("--cwd", default="", help="Command working directory relative to repo root.")
	parser.add_argument("command", nargs=argparse.REMAINDER, help="Command after --, or raw command tokens.")
	return parser.parse_args()


def get_command_from_args(args: argparse.Namespace) -> str | None:
	if not args.command:
		return None
	command = args.command
	if command[0] == "--":
		command = command[1:]
	if not command:
		return None
	return " ".join(shlex.quote(part) for part in command)


def command_was_from_stdin(args: argparse.Namespace) -> bool:
	return not args.command


def get_command_from_stdin() -> str | None:
	if sys.stdin.isatty():
		return None

	raw = sys.stdin.read()
	if not raw.strip():
		return None

	try:
		payload = json.loads(raw)
	except json.JSONDecodeError:
		return raw.strip()

	return find_command(payload)


def find_command(payload: Any) -> str | None:
	if isinstance(payload, str):
		return payload
	if isinstance(payload, list):
		return " ".join(shlex.quote(str(part)) for part in payload)
	if not isinstance(payload, dict):
		return None

	for key in ("command", "cmd", "shell", "input"):
		val = payload.get(key)
		if isinstance(val, str) and val.strip():
			return val
		if isinstance(val, list):
			return " ".join(shlex.quote(str(part)) for part in val)

	for key in ("tool_input", "params", "arguments", "data"):
		command = find_command(payload.get(key))
		if command:
			return command

	return None


def check_command(command: str, *, cwd: str = "") -> dict[str, Any]:
	normalized = normalize_command(command)
	if normalized is None:
		return {"allowed": False, "rule": "", "command": command, "why": "unable to parse command"}

	argv, command_cwd = normalized
	effective_cwd = normalize_cwd(cwd or command_cwd)

	for rule in load_rules():
		if rule_matches(rule, argv, effective_cwd):
			return {
				"allowed": True,
				"rule": rule["id"],
				"command": command,
				"argv": argv,
				"cwd": effective_cwd,
				"why": rule.get("why", ""),
			}

	return {"allowed": False, "rule": "", "command": command, "argv": argv, "cwd": effective_cwd}


def load_rules() -> list[dict[str, Any]]:
	with RULES_PATH.open() as f:
		data = json.load(f)
	rules = data.get("rules")
	if not isinstance(rules, list):
		raise ValueError("allow-command-rules.json must contain a rules list")
	return rules


def normalize_command(command: str) -> tuple[list[str], str] | None:
	command = command.strip()
	if not command:
		return None

	try:
		argv = shlex.split(command)
	except ValueError:
		return None

	if len(argv) >= 3 and Path(argv[0]).name in {"bash", "sh"} and argv[1] in {"-lc", "-c"}:
		return normalize_command(argv[2])

	if len(argv) >= 3 and argv[0] == "cd":
		cwd = argv[1]
		rest = argv[2:]
		if rest and rest[0] in {"&&", ";"}:
			rest = rest[1:]
		if rest:
			return rest, cwd

	return argv, ""


def normalize_cwd(cwd: str) -> str:
	if not cwd or cwd == ".":
		return ""
	path = Path(cwd)
	if path.is_absolute():
		try:
			return str(path.resolve().relative_to(REPO_ROOT))
		except ValueError:
			return str(path)
	return str(path)


def rule_matches(rule: dict[str, Any], argv: list[str], cwd: str) -> bool:
	rule_cwd = normalize_cwd(str(rule.get("cwd", "")))
	if rule_cwd and rule_cwd != cwd:
		return False

	kind = rule.get("kind")
	if kind == "exact":
		return argv == rule.get("argv")
	if kind == "prefix":
		prefix = rule.get("argv")
		return isinstance(prefix, list) and argv[:len(prefix)] == prefix
	if kind == "regex":
		pattern = rule.get("pattern")
		return isinstance(pattern, str) and re.match(pattern, " ".join(shlex.quote(part) for part in argv)) is not None

	return False


if __name__ == "__main__":
	raise SystemExit(main())
