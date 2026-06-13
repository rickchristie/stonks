# Playtest

Playwright integration tests live here. The canonical writing standards live in `skills/playtest/SKILL.md` and are symlinked for AI CLIs that support skill folders.

Start the configured backend and web servers before running tests:

```bash
setsid -f ./scripts/dev-backend.sh
setsid -f ./scripts/dev-web.sh
cd web/playtest && npx playwright test
```

All specs should import `{ test, expect }` from `../fixtures` or `./fixtures` so unexpected browser console errors fail the test.
