export type PrdDocumentKind = 'feature' | 'decision';

export type PrdDocument = {
	slug: string;
	title: string;
	summary: string;
	body: string[];
};

export type PrdStaticDocument = PrdDocument & {
	kind: PrdDocumentKind;
	path: string;
	sourcePath: string;
};

export type PrdIndexItem = {
	title: string;
	href: string;
	summary: string;
};

export type PrdIndexGroup = {
	title: string;
	items: PrdIndexItem[];
};
