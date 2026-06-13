import type { PrdStaticDocument } from '../../types';

const document = {
	kind: 'feature',
	slug: 'local-runtime',
	path: '/prd/features/local-runtime',
	sourcePath: 'web/src/lib/prd/documents/features/local-runtime.ts',
	title: 'Local Runtime',
	summary: 'Every checkout owns its local ports and database through `.dev`.',
	body: [
		'All scripts source `scripts/dev-env.sh`.',
		'The loader fails if `.dev` is missing or incomplete.',
		'`./scripts/init-db.sh` resets and seeds the configured local database.',
		'Dev backend, web, Playwright, and psql helpers all use the same config.',
		'The browser PRD is available at `/prd` in development and includes `/prd/storybook` scenarios.'
	]
} satisfies PrdStaticDocument;

export default document;
