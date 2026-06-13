import { test as base, expect } from '@playwright/test';

type ConsoleErrorFixture = {
	consoleErrors: {
		expectErrors(): void;
	};
};

const BROWSER_NOISE_RE = /^Failed to load resource:|^\[vite\]/;

export const test = base.extend<ConsoleErrorFixture>({
	consoleErrors: [
		async ({ page }, use) => {
			const errors: string[] = [];
			page.on('console', (msg) => {
				if (msg.type() === 'error' && !BROWSER_NOISE_RE.test(msg.text())) {
					errors.push(msg.text());
				}
			});

			let expected = false;
			await use({
				expectErrors() {
					expected = true;
				}
			});

			if (!expected) {
				expect(errors, `Unexpected console errors:\n${errors.join('\n')}`).toEqual([]);
			}
		},
		{ auto: true }
	]
});

export { expect } from '@playwright/test';
