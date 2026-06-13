# Production Release Guide

This directory keeps production context in the repository so agents can deploy, troubleshoot, and evolve infrastructure with the same information.

The scripts are intentionally generic. Rename service names, SSH host, domain, and paths before using them for a real app.

## Architecture

```text
User -> reverse proxy -> static web files
                    -> Go backend /api/*
Go backend -> PostgreSQL
```

## Required Customization

- `APP_PROD_SSH_HOST`: production SSH host. Defaults to `app-template-prod-1`.
- `APP_PROD_ROOT`: production app path. Defaults to `/opt/app-template`.
- Systemd service name: defaults to `app-template-backend`.
- Nginx/reverse proxy config: add one before production use.
- Secrets: keep production `.env` and database credentials on the server, not in this repo.

## Release Flow

1. Write a release note or migration plan for changes that require production sequencing.
2. Add SQL migrations under `release/migrations/`.
3. Run local verification:

```bash
cd backend && go test ./...
cd web && npm run check && npm test
cd web/playtest && npx playwright test
```

4. Deploy with `./scripts/release.sh`.
5. Run `./scripts/prod-smoke.sh --expect-version X.Y.Z`.
6. Append any server-level changes to `release/app-template-prod-1.changelog.md`.

`prod-smoke.sh` scans recent journal logs for both JSON levels like `"level":"error"` and the template logger's key-value levels like `level=error logger=...`. Keep `scripts/prod-smoke-test.sh` updated when changing these failure signatures.

## Production Query Helper

Use `scripts/prod-psql.sh` for production SQL. It sources credentials on the server and does not print them locally.

```bash
./scripts/prod-psql.sh -c "SELECT 1;"
./scripts/prod-psql.sh --file /tmp/query.sql
```

## Migration Rules

- Migration files are append-only and date-prefixed.
- Separate DML cleanup from destructive DDL when the old backend cannot tolerate the new schema.
- Stop the backend before destructive DDL when required.
- Record exact commands and outcomes in the changelog.
