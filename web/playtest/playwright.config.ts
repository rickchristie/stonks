import { defineConfig } from '@playwright/test';
import { baseURL } from './env';

const includeManual = process.env.PLAYWRIGHT_INCLUDE_MANUAL === '1';

export default defineConfig({
	testDir: '.',
	testIgnore: includeManual ? [] : ['**/manual/**'],
	timeout: 30000,
	retries: 1,
	workers: 4,
	use: {
		baseURL,
		headless: true,
		channel: 'chromium',
		screenshot: 'only-on-failure',
		trace: 'retain-on-failure'
	},
	webServer: undefined
});
