import type { PrdStaticDocument } from '../../types';

const document = {
	kind: 'decision',
	slug: '001-required-dev-config',
	path: '/prd/decisions/001-required-dev-config',
	sourcePath: 'web/src/lib/prd/documents/decisions/001-required-dev-config.ts',
	title: '001. Require `.dev`',
	summary: 'The template requires explicit local runtime configuration.',
	body: [
		'Implicit fallback hides configuration mistakes and makes parallel checkouts easier to collide.',
		'The template still provides `.dev.example`, but scripts read only `.dev`.',
		'New apps should keep `.dev` local to each checkout and update ports plus database names immediately after copying.'
	]
} satisfies PrdStaticDocument;

export default document;
