# Extracted From Mark

This document records the reusable practices found in `~/Personal/mark/*` and how this template applies them. It is intentionally product-neutral; Mark-specific data, OPML content, logs, generated builds, screenshots, backups, and private deployment details were not copied.

## Repository Shape

- Keep backend and frontend in one repo so agents can change API contracts, UI, tests, docs, and release notes together.
- Keep root scripts in `scripts/`; agents should not invent one-off commands when a repo-owned script can capture logging and local config.
- Keep `release/` in the repo so infrastructure notes, migrations, changelogs, and troubleshooting context travel with the app.
- Keep root `AGENTS.md` detailed and opinionated; it is the session-to-session memory that prevents quality drift.
- Keep `.dev` and `.dev.example` at the repo root. Runtime differences between parallel checkouts belong in `.dev`, not in scripts. Unlike Mark, this template deliberately fails when `.dev` is missing so new apps make their local runtime explicit.
- Keep a tool-neutral allow-list script so Codex, opencode, and other agents can share narrow command approval rules. Codex uses `.codex/config.toml`; OpenCode uses `opencode.json` for `/tmp` verification artifacts.

## Backend Practices

- Use Go and PostgreSQL as the foundation.
- Split code by responsibility:
  - `data`: internal DB-shaped entities and enums.
  - `accessor`: transaction-oriented storage code, SQL scanning, and domain-local database access.
  - `service`: business workflows and validation.
  - `pservice`: public HTTP API, auth, route wiring, request/response entities.
  - `lib`: reusable technical helpers.
- Keep schema under `backend/accessor/db/pg/setup/schema/*.init.sql`, one file per domain.
- Keep local DB reset and seed scripts under `scripts/`; schema SQL stays under `backend/accessor/db/pg/setup/schema/`.
- Keep business rules in service code, not DB triggers or stored procedures.
- Accessor methods return full objects, not partial structs. Partial/derived views get explicit `data/*_derived.go` files.
- `Get*` accessors return an error when not found; never return `(nil, nil)`.
- Write accessors use explicit parameters rather than accepting large data structs.
- Service methods use `Method(ctx context.Context, in MethodIn) MethodOut`; inputs include `Trace *tr.Trace`, outputs include `Success bool`.
- Public API entities mirror backend structs but omit sensitive/internal fields.
- HTTP APIs use `/api/*`, POST for app actions, and a response body with an `error` field. The template starts with open endpoints; real apps should add auth before sensitive data.

## Backend Testing Practices

- Include pgflock-backed database integration tests for accessors and services.
- Each test case gets fresh database state through `pgtest.WithDb`.
- Each domain owns `state_test.go` and `test_util.go` helpers.
- Test helper row structs use `Idx` plus required fields and defaults for optional fields.
- Assert database state first, service output second, and mocks/call counts third.
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

## Not Copied From Mark

- Mark product code beyond a small generic `note` example.
- Personal OPML data, screenshots, generated builds, test artifacts, logs, backups, and local secrets.
- Mark-specific Cloudflare R2 file flows, collaborative editor code, admin app, and private SSH host names.
- Mark's large product-specific Codex PermissionRequest hook. The template keeps a smaller generic hook pointed at `scripts/allow-command.py`.

## Future Extraction Candidates

- A template command that renames module/package/product identifiers safely.
- Optional auth module with HttpOnly cookie JWT, refresh tokens, and rate-limited login.
- Optional admin module using the reusable filter/get API pattern.
- Optional object storage module for signed uploads.
- A richer project-local Codex hook safety model if generic read-only command approvals are widened.
