#!/usr/bin/env python3
import json
import subprocess
import sys
import unittest
from pathlib import Path


REPO_ROOT = Path(__file__).resolve().parents[1]
SCRIPT = REPO_ROOT / "scripts" / "allow-command.py"


def run_allow(*args: str, stdin: str = "") -> subprocess.CompletedProcess[str]:
	return subprocess.run(
		[sys.executable, str(SCRIPT), *args],
		input=stdin,
		text=True,
		stdout=subprocess.PIPE,
		stderr=subprocess.PIPE,
		cwd=REPO_ROOT,
		check=False,
	)


class AllowCommandTest(unittest.TestCase):
	def test_allows_repo_script_prefix(self) -> None:
		result = run_allow("--json", "--", "./scripts/dev-health.sh")
		self.assertEqual(result.returncode, 0, result.stderr)
		payload = json.loads(result.stdout)
		self.assertTrue(payload["allowed"])
		self.assertEqual(payload["rule"], "dev-health")

	def test_allows_background_dev_server_wrappers(self) -> None:
		cases = [
			("dev-backend-background", "setsid -f ./scripts/dev-backend.sh"),
			("dev-web-background", "setsid -f ./scripts/dev-web.sh"),
		]

		for expected_rule, command in cases:
			with self.subTest(command=command):
				result = run_allow("--json", "--", "bash", "-lc", command)
				self.assertEqual(result.returncode, 0, result.stderr)
				payload = json.loads(result.stdout)
				self.assertTrue(payload["allowed"])
				self.assertEqual(payload["rule"], expected_rule)

	def test_allows_setup_wrapper(self) -> None:
		result = run_allow("--json", "--", "./scripts/setup.sh")
		self.assertEqual(result.returncode, 0, result.stderr)
		payload = json.loads(result.stdout)
		self.assertTrue(payload["allowed"])
		self.assertEqual(payload["rule"], "setup")

	def test_allows_local_db_reset_wrapper(self) -> None:
		result = run_allow("--json", "--", "./scripts/init-db.sh")
		self.assertEqual(result.returncode, 0, result.stderr)
		payload = json.loads(result.stdout)
		self.assertTrue(payload["allowed"])
		self.assertEqual(payload["rule"], "init-db")

	def test_allows_bounded_listener_and_process_probes(self) -> None:
		cases = [
			("lsof-listener", "lsof -nP -iTCP:8800 -sTCP:LISTEN"),
			("ps-dev-process-tree", "ps -o pid,ppid,pgid,cmd -p 123,4567"),
		]

		for expected_rule, command in cases:
			with self.subTest(command=command):
				result = run_allow("--json", "--", "bash", "-lc", command)
				self.assertEqual(result.returncode, 0, result.stderr)
				payload = json.loads(result.stdout)
				self.assertTrue(payload["allowed"])
				self.assertEqual(payload["rule"], expected_rule)

	def test_allows_cd_backend_go_test(self) -> None:
		result = run_allow("--json", "--", "bash", "-lc", "cd backend && go test ./...")
		self.assertEqual(result.returncode, 0, result.stderr)
		payload = json.loads(result.stdout)
		self.assertTrue(payload["allowed"])
		self.assertEqual(payload["cwd"], "backend")
		self.assertEqual(payload["rule"], "backend-tests")

	def test_denies_unknown_command(self) -> None:
		result = run_allow("--json", "--", "rm", "-rf", "/tmp/example")
		self.assertEqual(result.returncode, 1)
		payload = json.loads(result.stdout)
		self.assertFalse(payload["allowed"])

	def test_denies_env_file_read(self) -> None:
		result = run_allow("--json", "--", "cat", ".env")
		self.assertEqual(result.returncode, 1)
		payload = json.loads(result.stdout)
		self.assertFalse(payload["allowed"])

	def test_denies_ssh_key_read(self) -> None:
		result = run_allow("--json", "--", "cat", "~/.ssh/id_rsa")
		self.assertEqual(result.returncode, 1)
		payload = json.loads(result.stdout)
		self.assertFalse(payload["allowed"])

	def test_denies_destructive_find(self) -> None:
		result = run_allow("--json", "--", "find", "web/src", "-type", "f", "-delete")
		self.assertEqual(result.returncode, 1)
		payload = json.loads(result.stdout)
		self.assertFalse(payload["allowed"])

	def test_denies_arbitrary_kill(self) -> None:
		result = run_allow("--json", "--", "kill", "-TERM", "-123")
		self.assertEqual(result.returncode, 1)
		payload = json.loads(result.stdout)
		self.assertFalse(payload["allowed"])

	def test_codex_payload_shape(self) -> None:
		payload = json.dumps({"tool_input": {"command": "./scripts/dev-health.sh"}})
		result = run_allow("--codex", stdin=payload)
		self.assertEqual(result.returncode, 0, result.stderr)
		parsed = json.loads(result.stdout)
		decision = parsed["hookSpecificOutput"]["decision"]["behavior"]
		self.assertEqual(decision, "allow")

	def test_opencode_like_payload_shape(self) -> None:
		payload = json.dumps({"params": {"cmd": "cd web/playtest && npx playwright test prd/"}})
		result = run_allow("--json", stdin=payload)
		self.assertEqual(result.returncode, 0, result.stderr)
		parsed = json.loads(result.stdout)
		self.assertTrue(parsed["allowed"])
		self.assertEqual(parsed["rule"], "playwright-tests")


if __name__ == "__main__":
	unittest.main()
