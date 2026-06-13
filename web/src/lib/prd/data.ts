import { prdStaticDocuments } from './static-documents';
import type { PrdIndexGroup, PrdIndexItem, PrdStaticDocument } from './types';

export const featureDocs = prdStaticDocuments.filter((doc) => doc.kind === 'feature');
export const decisionDocs = prdStaticDocuments.filter((doc) => doc.kind === 'decision');

export const prdDocuments = [...featureDocs, ...decisionDocs];
export const prdRoutePaths = new Set(['/prd', '/prd/storybook', ...prdDocuments.map((doc) => doc.path)]);

function toIndexItem(doc: PrdStaticDocument): PrdIndexItem {
	return { title: doc.title, href: doc.path, summary: doc.summary };
}

export const prdIndexGroups: PrdIndexGroup[] = [
	{ title: 'Features', items: featureDocs.map(toIndexItem) },
	{ title: 'Decisions', items: decisionDocs.map(toIndexItem) }
];

export function getFeatureDoc(slug: string): PrdStaticDocument | undefined {
	return featureDocs.find((doc) => doc.slug === slug);
}

export function getDecisionDoc(slug: string): PrdStaticDocument | undefined {
	return decisionDocs.find((doc) => doc.slug === slug);
}
