import { test, expect } from '../fixtures';

test('renders Hello World from backend and database', async ({ page }) => {
	await page.goto('/app');

	await expect(page.getByTestId('hello-message')).toHaveText('Hello, World!');
	await expect(page.getByTestId('database-message')).toHaveText('This message came from PostgreSQL through the Go backend.');
	await expect(page.getByTestId('note-count')).toHaveText('1');
});
