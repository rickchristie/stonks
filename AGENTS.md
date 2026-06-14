# Critical Behavior
- **ALWAYS fully finish your tasks** when executing anything. Never stop to ask "would you like to continue?" or anything similar.
  You are given tasks, fully complete them, don't waste our time. Exception: Destructive actions with irreversible consequences.
- **ALWAYS verify your work and assumptions**, don't just read the code, actually test what you're doing.
  Write tests, write scripts to test behaviors, run playwright to check console, take screenshots, run the commands.
- **NEVER blame without evidence**, don't say something like "X fail due to Y", find evidence!
  **Always investigate to find root cause**, if unable to find evidence, state why and clarify it's a hypothesis.
- **ALWAYS write proper documentation**, write *why* it was done this way, and *how* only if it's not obvious.
  Write for a human/yourself when revisiting this code in the future, what's important so futureself work faster, fewer mistakes, less tokens?
- **NEVER execute staging, unstaging, or stashing git changes**, this will cause **loss of verification work** in our collaboration!
- **NEVER remove or edit changes that are not yours**, you are not the only one working in the repository.
  When you absolutely require to touch changes that are not yours, ASK permission.
- **NEVER delete valid comments**, contextual comments are important for maintainability.
  Removing comments that contains business context (the "why") will cause **loss of context**.
  Only remove comments if they are no longer valid.
- **ALWAYS use simple, intent-revealing names and comments**, e.g. use `shouldOrderInsertTimeline` instead of `orderQualifiesForLifecycleTimeline`.
  Prefer concise, immediately understood, simple names that describe a concrete action, decision, noun, business purpose.
  Avoid abstract, thesaurus-style, Java-like word salad. Applies to variables, functions, methods, tests, and comments.
- **ALWAYS prioritize readability, maintainability, simplicity, elegance** of your code.
  Aim for low-cyclomatic complexity in your code, exit early whenever you can, it's okay to repeat lines if we reduce cyclomatic
  complexity or prevent interleaving conditionals.
- **DO NOT excessively nil check**. Treat pointer/interface parameters/struct fields as required unless explicitly documented optional.
  If caller sends nil, then nil pointer panic is fine. Fail fast, unless we truly need the recovery.

## Stonks Goal

This repository is for Stonks, Rick's app for researching stocks, managing portfolios, and later automating trades.

## Documentation

- Keep root docs concise, terse, and complete.
- Product requirements live in `web/src/lib/prd` and `web/src/routes/prd`; do not create a detached root `/prd` folder.
- Architecture and production decisions belong in PRD decision docs or `release/`, not only in chat.
- Routeable PRD document bodies live under `web/src/lib/prd/documents` and must be registered in `web/src/lib/prd/static-documents.ts`; keep `web/src/lib/prd/data.test.ts` passing when adding or moving docs.
- Update docs whenever scripts, test setup, release flow, or architecture contracts change.

## Backend

- Location: `backend/*`.
- Architecture:
  - `backend/data`: internal DB-shaped entities and enums.
  - `backend/accessor`: transaction-oriented database access.
  - `backend/service`: business workflows and validation.
  - `backend/pservice`: public API routes, request/response entities, middleware, entrypoint.
  - `backend/accessor/db/pg/setup/schema`: domain SQL init files.

Standards:
- Check every returned error. Never swallow errors.
- Keep business logic in Go service code, not DB triggers/procedures.
- Use `time.Time` for timestamps and field names `CreatedTs`, `LastUpdatedTs`.
- Use zero values for optional primitive data unless optionality is explicit.
- Read-only transactions should be rolled back, not committed.
- Accessors return full objects. Derived structs belong in `backend/data/*_derived.go`.
- `Get*` accessors return `err != nil` when not found.
- Write accessors use explicit parameters rather than accepting entity structs.
- Service methods use `Method(ctx context.Context, in MethodIn) MethodOut`.
- Service input structs include `Trace *tr.Trace`; output structs include `Success bool`.
- Public service handlers depend on consumer-facing service interfaces from `backend/service/*/lib`, not implementation packages.
- Service code logs infrastructure failures with trace context before returning internal errors.
- Public API routes use `/api/*`, POST for app actions, and responses with an `error` field.
- Public API entities mirror backend data names where practical.

Testing:
- Accessor and service behavior must have tests.
- Database tests use pgflock through `backend/accessor/db/pg/pgtest`.
- Each test case should get fresh database state.
- Each domain owns `state_test.go` and `test_util.go`.
- Assert database state first, service output second, mock calls third.
- Assert all relevant fields, not only happy-path fragments.

Import aliases:
- `accs[Domain]` for accessor imports.
- `svc[Domain]` for service imports.
- `ps[Domain]` for public service imports.

## Web App

- Location: `web/*`.
- The frontend must compile to static HTML/CSS/JS.
- Root SvelteKit routes explicitly disable SSR/prerender for the static SPA contract, and `hooks.client.ts` owns dynamic-import recovery for route chunk failures.
- Use Svelte 5 runes: `$props()`, `$state()`, `$derived()`, `$effect()`, `$bindable()`.
- Use `web/src/lib/api/client.ts` for all fetch logic. Transport failures are separate from app-level response `error` values.
- Keep API types under `web/src/lib/types` or next to API modules when the app is small.
- Keep reusable UI in `web/src/lib/components`.
- Keep tests next to source as `*.test.ts`.
- Use Vitest for pure/frontend unit tests and Playwright for browser behavior.
- Playwright specs import `{ test, expect }` from `web/playtest/fixtures.ts` so unexpected console errors fail tests.
- Playtest standards live in `skills/playtest/SKILL.md`; use that skill when writing, reviewing, debugging, or running Playwright specs.
- PRD and storybook routes are dev-only.

## Development

- Prefer direct commands; use `bash -lc` only when shell semantics are required.
- Repo-local config is read from `.dev`. Scripts must fail if `.dev` is missing or incomplete.
- Keep `.dev` committed in this repo with non-secret local defaults; real app secrets still belong in ignored `.env` files or server config.
- Start dev servers as background processes through scripts so logs and ports are consistent:
  - `setsid -f ./scripts/dev-backend.sh`
  - `setsid -f ./scripts/dev-web.sh`
- Run `./scripts/dev-health.sh` before and after server work.
- Query local DB with `./scripts/dev-psql.sh`.
- Reset and seed local DB with `./scripts/init-db.sh`.
- Shared agent command approvals live in `scripts/allow-command-rules.json`; update `scripts/allow-command-test.py` when adding a rule.

Verification checklist:
- Backend: `cd backend && go test ./...`.
- Web unit/static checks: `cd web && npm run check && npm test`.
- Browser integration: `cd web/playtest && npx playwright test`.
- Health and logs: `./scripts/dev-health.sh`.

## Release

- Production context lives in `release/`.
- Migration scripts live in `release/migrations/`.
- Any production migration or manual server work needs release notes or changelog context.
- Production helper scripts must not print credentials.
- See `release/README.md` before changing deployment behavior.
