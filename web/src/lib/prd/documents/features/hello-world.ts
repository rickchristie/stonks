import type { PrdStaticDocument } from '../../types';

const document = {
	kind: 'feature',
	slug: 'hello-world',
	path: '/prd/features/hello-world',
	sourcePath: 'web/src/lib/prd/documents/features/hello-world.ts',
	title: 'Hello World Template',
	summary: 'The first screen proves frontend, backend, and PostgreSQL are connected.',
	body: [
		'The app route `/app` renders a Hello World surface.',
		'The page calls `/api/hello` through the shared API client.',
		'The backend route reads active rows from PostgreSQL before responding.',
		'The seeded row body appears as database proof in the UI.'
	]
} satisfies PrdStaticDocument;

export default document;
