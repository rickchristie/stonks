import { cleanup, render, screen } from '@testing-library/svelte';
import { afterEach, describe, expect, it, vi } from 'vitest';
import Page from './+page.svelte';

const appEnvironment = vi.hoisted(() => ({
	dev: true
}));

const helloApi = vi.hoisted(() => ({
	getHello: vi.fn()
}));

vi.mock('$app/environment', () => ({
	browser: true,
	building: false,
	get dev() {
		return appEnvironment.dev;
	},
	version: 'test'
}));

vi.mock('$lib/api/hello', () => ({
	getHello: helloApi.getHello
}));

describe('/app shell', () => {
	afterEach(() => {
		cleanup();
		vi.clearAllMocks();
	});

	it('shows PRD links in development', async () => {
		await renderAppPage(true);

		expect(screen.getByRole('link', { name: 'PRD' }).getAttribute('href')).toBe('/prd');
		expect(screen.getByRole('link', { name: 'Storybook' }).getAttribute('href')).toBe('/prd/storybook');
	});

	it('hides PRD links in production', async () => {
		await renderAppPage(false);

		expect(screen.queryByRole('navigation', { name: 'Stonks docs' })).toBeNull();
		expect(screen.queryByRole('link', { name: 'PRD' })).toBeNull();
		expect(screen.queryByRole('link', { name: 'Storybook' })).toBeNull();
	});
});

function renderAppPage(dev: boolean) {
	appEnvironment.dev = dev;
	helloApi.getHello.mockResolvedValue({
		error: '',
		message: 'Hello from test',
		databaseMessage: 'Seeded database row',
		noteCount: 1
	});
	render(Page);
}
