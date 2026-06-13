import helloWorldDocument from './documents/features/hello-world';
import localRuntimeDocument from './documents/features/local-runtime';
import requiredDevConfigDocument from './documents/decisions/001-required-dev-config';
import type { PrdStaticDocument } from './types';

export type { PrdStaticDocument } from './types';

// Keep body content in ./documents so content reviews stay file-local as PRD grows.
export const prdStaticDocuments = [
	helloWorldDocument,
	localRuntimeDocument,
	requiredDevConfigDocument
] satisfies PrdStaticDocument[];

export const prdStaticDocumentByPath = new Map(
	prdStaticDocuments.map((document) => [document.path, document])
);

export function getPrdStaticDocument(pathname: string): PrdStaticDocument | undefined {
	return prdStaticDocumentByPath.get(pathname);
}
