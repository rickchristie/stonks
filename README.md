# Stonks

Stonks is Rick's app for researching stocks, managing portfolios, and later automating trades. This repository is currently in preparation mode: the product has not been built yet, but the full-stack operating model is ready for step-by-step implementation.

## What Is Ready

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

See [docs/stonks-repo-prep.md](docs/stonks-repo-prep.md) for the preparation notes and the choices intentionally left for the real Stonks build.

## Layout

```text
backend/          Go backend: data, accessor, service, pservice
web/              Svelte 5 static frontend, PRD, storybook, unit tests
web/playtest/     Playwright integration tests
scripts/          Local dev, health, database, and test helpers
skills/           Canonical AI skills used by this repo
release/          Production deployment docs, changelog, migrations
docs/             Stonks preparation notes
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

The app first screen is still a Hello World verification page that calls `/api/hello`; that backend route reads the seeded database row before responding. Keep it until the first real Stonks workflow has equivalent backend, frontend, and browser coverage. Backend health is available at `/health-check` and `/version`. Frontend PRD is dev-only at `/prd`.

## Agent Allow List

`scripts/allow-command.py` is a tool-neutral allow-list checker for Codex, opencode, or any agent runner that supports command approval hooks. Rules live in `scripts/allow-command-rules.json` and tests live in `scripts/allow-command-test.py`.

Codex hook wiring lives in `.codex/config.toml`. OpenCode external `/tmp` artifact access lives in `opencode.json`.

Examples:

```bash
python3 scripts/allow-command.py --check -- ./scripts/dev-health.sh
python3 scripts/allow-command.py --json -- bash -lc 'cd backend && go test ./...'
```

## Stonks Build Notes

The repository identity has been renamed to Stonks, but the example `note` domain remains as working scaffolding. Do not remove it until the replacement Stonks domain has accessor, service, pservice, unit, and playtest coverage. It is intentionally small but demonstrates the architecture future product work should follow.
