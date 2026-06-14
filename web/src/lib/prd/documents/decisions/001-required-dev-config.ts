import type { PrdStaticDocument } from '../../types';

const document = {
	kind: 'decision',
	slug: '001-required-dev-config',
	path: '/prd/decisions/001-required-dev-config',
	sourcePath: 'web/src/lib/prd/documents/decisions/001-required-dev-config.ts',
	title: '001. Require `.dev`',
	summary: 'Stonks requires explicit local runtime configuration.',
	body: [
		'Implicit fallback hides configuration mistakes and makes parallel checkouts easier to collide.',
		'Stonks still provides `.dev.example`, but scripts read only `.dev`.',
		'Each checkout should keep `.dev` local and update ports plus database names before starting servers.'
	]
} satisfies PrdStaticDocument;

export default document;
