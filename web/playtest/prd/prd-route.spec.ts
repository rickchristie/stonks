import { test, expect } from '../fixtures';

test('renders PRD and storybook routes', async ({ page }) => {
	await page.goto('/prd');
	await expect(page.getByRole('heading', { name: 'Stonks' })).toBeVisible();
	await expect(
		page.getByLabel('PRD navigation').getByRole('link', { name: 'Stonks Hello World' })
	).toBeVisible();

	await page.goto('/prd/storybook');
	await expect(page.getByRole('heading', { name: 'Hello World States' })).toBeVisible();
	await expect(page.getByText('This message came from PostgreSQL through the Go backend.')).toBeVisible();
});
