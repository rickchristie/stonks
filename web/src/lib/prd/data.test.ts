import { describe, expect, it } from 'vitest';
import { prdDocuments, prdIndexGroups, prdRoutePaths } from './data';
import { prdStaticDocuments } from './static-documents';

const prdDocumentModules = import.meta.glob('./documents/**/*.ts');
const PRD_PATH_RE = /\/prd(?:\/[A-Za-z0-9._-]+)*/g;

function duplicateValues(values: string[]): string[] {
	const seen = new Set<string>();
	const duplicates = new Set<string>();
	values.forEach((value) => {
		if (seen.has(value)) duplicates.add(value);
		seen.add(value);
	});
	return [...duplicates].sort();
}

function prdLinksFromBody(body: string[]): string[] {
	return body.flatMap((line) => line.match(PRD_PATH_RE) ?? []);
}

describe('PRD static document registry', () => {
	it('has unique document slugs, paths, and page routes', () => {
		expect(duplicateValues(prdStaticDocuments.map((doc) => doc.slug))).toEqual([]);
		expect(duplicateValues(prdStaticDocuments.map((doc) => doc.path))).toEqual([]);
		expect(prdRoutePaths.size).toBe(prdDocuments.length + 2);
	});

	it('keeps static document source paths file-local', () => {
		prdStaticDocuments.forEach((document) => {
			const modulePath = `./${document.sourcePath.replace('web/src/lib/prd/', '')}`;

			expect(document.sourcePath).not.toContain('static-documents.ts#');
			expect(document.sourcePath.startsWith('web/src/lib/prd/documents/')).toBe(true);
			expect(prdDocumentModules[modulePath]).toBeDefined();
		});
	});

	it('registers every static document in the sidebar index', () => {
		const sidebarHrefs = new Set(prdIndexGroups.flatMap((group) => group.items.map((item) => item.href)));
		const missingPaths = prdStaticDocuments
			.map((document) => document.path)
			.filter((pathname) => !sidebarHrefs.has(pathname));

		expect(missingPaths).toEqual([]);
	});

	it('keeps PRD index links backed by registered routes', () => {
		const brokenItems = prdIndexGroups
			.flatMap((group) => group.items)
			.filter((item) => !prdRoutePaths.has(item.href))
			.map((item) => ({ title: item.title, href: item.href }));

		expect(brokenItems).toEqual([]);
	});

	it('keeps /prd links in document bodies backed by registered routes', () => {
		const brokenLinks = prdStaticDocuments.flatMap((document) =>
			prdLinksFromBody(document.body)
				.filter((pathname) => !prdRoutePaths.has(pathname))
				.map((pathname) => ({ document: document.path, href: pathname }))
		);

		expect(brokenLinks).toEqual([]);
	});
});
