# App Template

Reusable monorepo template for agentic product development. The template keeps backend, frontend, product docs, tests, local runtime config, and release context in one repository so every agent session has the same operating model.

## What This Extracts From Mark

- One repository for Go backend, Svelte frontend, Playwright, PRD, release docs, and infrastructure scripts.
- Repo-local `.dev` runtime config so multiple checkouts can run separate ports and databases.
- Go backend split by responsibility: `data`, `accessor`, `service`, `pservice`.
- PostgreSQL schema files per domain, with local reset scripts and pgflock-backed integration test helpers.
- Svelte static frontend with dev-only `/prd` and `/prd/storybook` routes.
- npm install hardening through `.npmrc` `min-release-age`.
- Scripts that write logs to `logs/`, probe health, run psql, start dev servers, and run Playwright.
- Shared command allow-list script for Codex, opencode, and other agent runners.
- Canonical `playtest` skill under `skills/playtest`, mirrored for OpenCode and Claude.
- Release directory for deployment docs, migration scripts, smoke checks, and production troubleshooting notes.
- Root `AGENTS.md` as the operating manual for future agent sessions.

See [docs/extracted-from-mark.md](docs/extracted-from-mark.md) for the full extraction checklist and the choices intentionally left as placeholders.

## Layout

```text
backend/          Go backend: data, accessor, service, pservice
web/              Svelte 5 static frontend, PRD, storybook, unit tests
web/playtest/     Playwright integration tests
scripts/          Local dev, health, database, and test helpers
skills/           Canonical AI skills used by this template
release/          Production deployment docs, changelog, migrations
docs/             Template rationale and extraction notes
.pgflock/         pgflock config for isolated PostgreSQL test DBs
```

## First Setup

1. Update `.dev` and adjust `DB_NAME`, `DEV_BACKEND_PORT`, and `DEV_WEB_PORT` if another checkout is already using the defaults. Scripts fail when `.dev` is missing or incomplete.
2. Install dependencies with `./scripts/setup.sh`.
3. Reset and seed local database with `./scripts/init-db.sh`.

The database reset script only allows local PostgreSQL targets because it uses the local `postgres` OS user.

## Development

Run backend and frontend through the repo scripts so they respect `.dev`:

```bash
./scripts/dev-health.sh
setsid -f ./scripts/dev-backend.sh
setsid -f ./scripts/dev-web.sh
./scripts/dev-health.sh
```

Useful checks:

```bash
cd backend && go test ./...
cd web && npm run check && npm test
cd web/playtest && npx playwright test
./scripts/dev-health.sh
```

The app first screen is a Hello World page that calls `/api/hello`; that backend route reads the seeded database row before responding. Backend health is available at `/health-check` and `/version`. Frontend PRD is dev-only at `/prd`.

## Agent Allow List

`scripts/allow-command.py` is a tool-neutral allow-list checker for Codex, opencode, or any agent runner that supports command approval hooks. Rules live in `scripts/allow-command-rules.json` and tests live in `scripts/allow-command-test.py`.

Codex hook wiring lives in `.codex/config.toml`. OpenCode external `/tmp` artifact access lives in `opencode.json`.

Examples:

```bash
python3 scripts/allow-command.py --check -- ./scripts/dev-health.sh
python3 scripts/allow-command.py --json -- bash -lc 'cd backend && go test ./...'
```

## Template Customization

Rename the module and product strings before using this as a real app:

- `backend/go.mod` module path.
- `web/package.json` package names.
- `.dev` and `.dev.example` database/user names and default ports if needed.
- `release/README.md`, `scripts/release.sh`, production helper scripts, systemd service name, SSH host, and production paths.
- PRD documents under `web/src/lib/prd/documents`, the registry in `web/src/lib/prd/static-documents.ts`, and route files under `web/src/routes/prd`.

Do not remove the example `note` domain until the replacement domain has accessor, service, pservice, unit, and playtest coverage. It is intentionally small but demonstrates the architecture.
