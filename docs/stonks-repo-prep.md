# Stonks Repository Prep

This document records why the repository is prepared this way before the real Stonks product work begins. Stonks is intended to help Rick research stocks, manage portfolios, and later automate trades, but those workflows should be added only when the build is explicitly started.

## Repository Shape

- Keep backend and frontend in one repo so agents can change API contracts, UI, tests, docs, and release notes together.
- Keep root scripts in `scripts/`; agents should not invent one-off commands when a repo-owned script can capture logging and local config.
- Keep `release/` in the repo so infrastructure notes, migrations, changelogs, and troubleshooting context travel with the app.
- Keep root `AGENTS.md` detailed and opinionated; it is the session-to-session memory that prevents quality drift.
- Keep `.dev` and `.dev.example` at the repo root. Runtime differences between parallel checkouts belong in `.dev`, not in scripts.
- Keep a tool-neutral allow-list script so Codex, opencode, and other agents can share narrow command approval rules. Codex uses `.codex/config.toml`; OpenCode uses `opencode.json` for `/tmp` verification artifacts.

## Preserved Architecture

- Go and PostgreSQL remain the backend foundation.
- Backend code stays split by responsibility:
  - `data`: internal DB-shaped entities and enums.
  - `accessor`: transaction-oriented storage code, SQL scanning, and domain-local database access.
  - `service`: business workflows and validation.
  - `pservice`: public HTTP API, auth, route wiring, request/response entities.
  - `lib`: reusable technical helpers.
- Schema belongs under `backend/accessor/db/pg/setup/schema/*.init.sql`, one file per domain.
- Local DB reset and seed scripts stay under `scripts/`; schema SQL stays under `backend/accessor/db/pg/setup/schema/`.
- Business rules belong in service code, not DB triggers or stored procedures.
- Accessor methods return full objects, not partial structs. Partial or derived views get explicit `data/*_derived.go` files.
- `Get*` accessors return an error when not found; never return `(nil, nil)`.
- Write accessors use explicit parameters rather than accepting large data structs.
- Service methods use `Method(ctx context.Context, in MethodIn) MethodOut`; inputs include `Trace *tr.Trace`, outputs include `Success bool`.
- Public API entities mirror backend structs but omit sensitive or internal fields.
- HTTP APIs use `/api/*`, POST for app actions, and a response body with an `error` field. Add auth before introducing sensitive financial data.

## Testing Practices

- Include pgflock-backed database integration tests for accessors and services.
- Each test case gets fresh database state through `pgtest.WithDb`.
- Each domain owns `state_test.go` and `test_util.go` helpers.
- Test helper row structs use `Idx` plus required fields and defaults for optional fields.
- Assert database state first, service output second, and mocks or call counts third.
- Assert complete output fields instead of partial success checks.
- Use descriptive specs such as `"creates note with title and body"`, not `"basic"`.

## Frontend Practices

- Use Svelte 5 runes and SvelteKit static adapter.
- Compile frontend to static HTML/CSS/JS so production hosting stays cheap and simple.
- Keep API fetch logic in `web/src/lib/api/client.ts`; domain modules call the client.
- Keep colocated unit tests in `web/src/**/*.test.ts`.
- Use Vitest with jsdom for pure frontend logic and Playwright for browser behavior.
- Keep Playwright tests in `web/playtest`, with `fixtures.ts` checking unexpected console errors.
- Harden npm installs with `.npmrc`: `engine-strict=true` and `min-release-age=7`.

## PRD, Storybook, And Decisions

- Host PRD in the frontend at `/prd`, but gate it to development only.
- Keep PRD source in TypeScript/Svelte files under `web/src/lib/prd/documents`, register routeable docs in `web/src/lib/prd/static-documents.ts`, and render through `web/src/routes/prd`; do not add an unconnected root `/prd` folder.
- Keep storybook-like scenarios under `/prd/storybook` so agents can inspect component states without backend wiring.
- Record architectural decisions as PRD decision documents, not scattered chat context.
- Keep PRD concise but complete enough for another agent to understand expected behavior and test coverage.

## Scripts And Local Runtime

- `scripts/dev-env.sh` is the single source of local config. It requires `.dev` and exports backend/web/database/test URLs.
- `scripts/allow-command.py` centralizes recurring approval rules for agent tools. New rules go into `scripts/allow-command-rules.json` and must be covered by `scripts/allow-command-test.py`.
- Dev server scripts kill only the configured port, not every process for the project type.
- Dev server logs go to `logs/backend-$PORT.log` and `logs/web-$PORT.log`.
- `scripts/dev-health.sh` collects bounded evidence: configured URLs, listeners, HTTP status, DB access, process probes, and log tails.
- `scripts/dev-psql.sh` connects using `.dev` without printing credentials.
- `scripts/test-playwright.sh` provides a guided runner and logs to `logs/`.
- `scripts/setup.sh` installs Go, web, and playtest dependencies without reading `.dev`; npm still enforces each `.npmrc`.

## Release Practices

- Keep production setup docs, migrations, and changelogs under `release/`; keep executable helpers under `scripts/`.
- Keep migrations in `release/migrations/` and write release notes when production requires manual sequencing.
- Provide smoke scripts for production health/version/service/log checks.
- Keep production credentials on the server; helpers should source them remotely and not print them locally.
- Keep changelogs under `release/` for server changes and production incidents.

## Deferred Product Work

- Keep the example `note` domain until the first real Stonks domain has equivalent accessor, service, pservice, unit, and browser coverage.
- Add auth before portfolio, brokerage, or trading data is persisted.
- Treat automated trading as a later workflow with explicit risk controls, audit logs, dry-run mode, and kill-switch behavior.
- Do not add equity data providers, brokerage integrations, or portfolio models until the build step requests them.
