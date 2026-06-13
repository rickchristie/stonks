import { test, expect } from '../fixtures';

test('reloads once when a dynamic route chunk fails', async ({ page }) => {
	await page.goto('/app');
	await expect(page.getByRole('heading', { name: 'Hello World' })).toBeVisible();

	const loadPromise = page.waitForEvent('load');
	await page.evaluate(() => {
		const event = new CustomEvent('vite:preloadError', {
			cancelable: true,
			detail: new Error('Failed to fetch dynamically imported module')
		});
		window.dispatchEvent(event);
	});
	await loadPromise;

	await expect(page).toHaveURL(/\/app$/);
	await expect(page.getByRole('heading', { name: 'Hello World' })).toBeVisible();
	await expect(page.getByRole('heading', { name: /Error 500/i })).toHaveCount(0);
});
