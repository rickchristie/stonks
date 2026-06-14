---
name: playtest
description: Use when writing, reviewing, debugging, or running Playwright integration tests. Applies to web/playtest specs, manual exploratory specs, PRD coverage mapping, console-error handling, screenshots/traces, and browser verification.
---

# Playtest

Use this skill for every Playwright integration-test task in this repository.

## Core Rules

- Start from `.dev`; never hardcode backend or web ports.
- Run dev servers as background processes through `./scripts/dev-backend.sh` and `./scripts/dev-web.sh`.
- Run `./scripts/dev-health.sh` before and after server work.
- Specs live under `web/playtest/` and run from that directory with `npx playwright test`.
- Manual throwaway specs live under `web/playtest/manual/` and are ignored unless `PLAYWRIGHT_INCLUDE_MANUAL=1` is set.
- All committed specs must import `{ test, expect }` from the repo fixture, not from `@playwright/test`.
- Use user-visible assertions first; use implementation selectors only when the UI has no stable accessible surface.
- Tests must be deterministic, isolated, and repeatable without assuming a clean browser profile.
- Every user-facing behavior in PRD should have matching playtest coverage or an explicit reason why unit tests are enough.

## Commands

Run from the repository root unless noted.

```bash
./scripts/dev-health.sh
setsid -f ./scripts/dev-backend.sh
setsid -f ./scripts/dev-web.sh
./scripts/dev-health.sh
cd web/playtest
npx playwright test
npx playwright test app/hello.spec.ts
npx playwright test --headed
npx playwright test --debug
PLAYWRIGHT_INCLUDE_MANUAL=1 npx playwright test manual/my-debug.spec.ts
```

If Playwright reports a missing browser revision, run the install command it prints from `web/playtest`.

## Imports And Console Errors

Every spec imports from the nearest relative path to `fixtures.ts`.

```typescript
import { test, expect } from '../fixtures';
import type { Page } from '@playwright/test';
```

The fixture collects browser `console.error` messages and fails the test at teardown. This catches real client failures even when the visible assertion passes.

For tests that intentionally trigger console errors, opt out explicitly:

```typescript
test('handles failed API response', async ({ page, consoleErrors }) => {
	consoleErrors.expectErrors();
	// Trigger and assert the expected failure UI.
});
```

Do not silence console errors globally. If a browser or dev-server message is noise, add a narrow ignore pattern to `fixtures.ts` with a comment explaining why it is not an app error.

## PRD Coverage

Each feature-level PRD entry should map to one or more specs. Keep this map updated when adding new feature suites.

| Test File | PRD Route | Coverage |
| --- | --- | --- |
| `app/hello.spec.ts` | `/prd/features/hello-world` | Hello World app loads, calls backend, and renders a PostgreSQL-backed response |
| `prd/prd-route.spec.ts` | `/prd` and `/prd/storybook` | PRD index, feature navigation, and storybook route render in development |

When adding a feature:

1. Add or update its PRD page.
2. Add a playtest spec for the primary user workflow.
3. Add edge-case specs for destructive, permission, persistence, sync, or cross-session behavior.
4. Update this coverage table.

## Idempotent Test Pattern

Tests must pass on repeated runs and when multiple agents run suites against separate `.dev` configs.

- Create records needed by the test instead of depending on incidental previous state.
- Use unique names with `Date.now()` or a module-level test-run ID.
- Clear app-owned browser persistence before tests that depend on default UI state.
- Avoid shared mutable module state unless it is a read-only test-run identifier.
- Assert the final persisted state when testing writes.
- Prefer APIs or helper functions to set up data when UI setup would make the test slow or brittle.

Example:

```typescript
const testRunId = Date.now();

test('creates a note from the app', async ({ page }) => {
	const title = `Playtest Note ${testRunId}`;
	await page.goto('/app');
	await page.getByRole('button', { name: 'New Note' }).click();
	await page.getByRole('textbox', { name: 'Title' }).fill(title);
	await page.getByRole('button', { name: 'Save' }).click();
	await expect(page.getByText(title, { exact: true })).toBeVisible();
});
```

## Destructive Test Isolation

Any test that confirms a destructive or irreversible action must create its own throwaway data first.

Destructive actions include:

- Delete, archive, suspend, deactivate, revoke, reset, force logout, and irreversible migrations.
- UI confirm buttons that execute those actions.
- Direct API calls that mutate existing shared records.

Safe shared-data operations:

- Read-only filter/get/list behavior.
- Creating new throwaway records.
- Opening a modal and clicking Cancel.

Rule: opening a destructive modal is safe; clicking the confirm button is destructive.

## Selector Priority

Use selectors in this order:

1. `getByRole()` for buttons, textboxes, links, headings, menus, and dialogs.
2. `getByLabel()` for named regions or form controls.
3. `getByText()` for visible copy.
4. `getByTestId()` when accessible selectors would be ambiguous.
5. CSS locators only for structural details that users cannot name.

Use `exact: true` for short or common names.

```typescript
// Bad: may match another item with a longer name.
await page.getByRole('button', { name: 'Save' }).click();

// Good.
await page.getByRole('button', { name: 'Save', exact: true }).click();
```

Scope locators when the same text appears in navigation and content.

```typescript
await expect(
	page.getByLabel('PRD navigation').getByRole('link', { name: 'Stonks Hello World' })
).toBeVisible();
```

## Waiting And Assertions

- Use Playwright auto-waiting assertions instead of fixed sleeps.
- Wait for URLs after navigation-triggering actions.
- Assert visible user outcomes before internal state.
- For backend-backed UI, assert both the user-visible UI and the API/database proof when practical.
- For async saves, wait for the explicit saved state or verify data after reload.

Examples:

```typescript
await page.getByRole('link', { name: 'PRD' }).click();
await page.waitForURL(/\/prd$/);
await expect(page.getByRole('heading', { name: 'Stonks' })).toBeVisible();

await expect(page.getByText('Saved')).toBeVisible({ timeout: 5000 });
```

## Screenshots, Traces, And Manual Probes

- Playwright config should keep `screenshot: 'only-on-failure'` and `trace: 'retain-on-failure'`.
- For layout-sensitive work, capture explicit screenshots into `.playwright-mcp/` or `web/playtest/.shots/`.
- Before redesigning an existing page, take a screenshot of the current design.
- Manual probe specs belong in `web/playtest/manual/`; do not commit scratch specs outside that folder.
- If manual evidence becomes part of regression coverage, promote the probe into a normal spec and remove manual-only assumptions.

## Debugging Workflow

When a playtest fails:

1. Read the failure, screenshot, and trace path.
2. Re-run the single spec, preferably with `--headed` for UI issues.
3. Check browser console failures from the fixture output.
4. Inspect backend and web logs under `logs/`.
5. Use temporary `console.log` only to gather evidence, then remove it.
6. Fix the root cause or tighten the test if the failure is a selector/test-isolation bug.
7. Re-run the affected spec and then the full `web/playtest` suite.

Do not explain a failure as flakiness without evidence.

## Review Checklist

- Specs import from `fixtures.ts`.
- No hardcoded ports or absolute local checkout paths.
- Tests are idempotent and have unique data where needed.
- Destructive behavior uses throwaway data.
- Selectors are accessible and scoped.
- PRD coverage table is updated.
- Manual scratch specs remain under `manual/`.
- `npx playwright test` passes, or the blocker is documented with exact evidence.
