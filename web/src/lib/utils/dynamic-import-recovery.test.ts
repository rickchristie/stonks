import { describe, expect, it, vi } from 'vitest';

import {
	extractDynamicImportFailureMessage,
	installDynamicImportRecovery,
	isDynamicImportFailure,
	shouldReloadForDynamicImportFailure,
	type StorageLike
} from './dynamic-import-recovery';

function createStorage(): StorageLike {
	const values = new Map<string, string>();
	return {
		getItem(key: string) {
			return values.get(key) ?? null;
		},
		setItem(key: string, value: string) {
			values.set(key, value);
		}
	};
}

describe('dynamic-import-recovery', () => {
	it('reads messages from strings, errors, and plain objects', () => {
		expect(extractDynamicImportFailureMessage('plain')).toBe('plain');
		expect(extractDynamicImportFailureMessage(new Error('boom'))).toBe('boom');
		expect(extractDynamicImportFailureMessage({ message: 'object boom' })).toBe('object boom');
		expect(extractDynamicImportFailureMessage({})).toBe('');
	});

	it('matches known dynamic import failure messages only', () => {
		expect(isDynamicImportFailure('TypeError: Failed to fetch dynamically imported module')).toBe(true);
		expect(isDynamicImportFailure(new Error('Importing a module script failed.'))).toBe(true);
		expect(isDynamicImportFailure({ message: 'error loading dynamically imported module' })).toBe(true);
		expect(isDynamicImportFailure('Network Error')).toBe(false);
	});

	it('reloads once per URL inside the recovery window', () => {
		const storage = createStorage();

		expect(
			shouldReloadForDynamicImportFailure(
				storage,
				'http://localhost:35173/prd/storybook',
				1000,
				'Failed to fetch dynamically imported module'
			)
		).toBe(true);

		expect(
			shouldReloadForDynamicImportFailure(
				storage,
				'http://localhost:35173/prd/storybook',
				5000,
				'Failed to fetch dynamically imported module'
			)
		).toBe(false);

		expect(
			shouldReloadForDynamicImportFailure(
				storage,
				'http://localhost:35173/app',
				6000,
				'Failed to fetch dynamically imported module'
			)
		).toBe(true);

		expect(
			shouldReloadForDynamicImportFailure(
				storage,
				'http://localhost:35173/prd/storybook',
				17001,
				'Failed to fetch dynamically imported module'
			)
		).toBe(true);
	});

	it('ignores non-dynamic-import failures', () => {
		expect(
			shouldReloadForDynamicImportFailure(createStorage(), 'http://localhost:35173/app', 1000, 'Network Error')
		).toBe(false);
	});

	it('reloads once for a preload error and prevents the default event', () => {
		const target = new EventTarget();
		const reload = vi.fn();
		const dispose = installDynamicImportRecovery({
			target,
			storage: createStorage(),
			getHref: () => 'http://localhost:35173/app',
			now: () => 1000,
			reload
		});

		const event = new CustomEvent('vite:preloadError', { cancelable: true });
		Object.defineProperty(event, 'payload', {
			value: new Error('Failed to fetch dynamically imported module')
		});

		target.dispatchEvent(event);
		target.dispatchEvent(event);

		expect(reload).toHaveBeenCalledTimes(1);
		expect(event.defaultPrevented).toBe(true);

		dispose();
	});

	it('handles window error and unhandled rejection events', () => {
		const target = new EventTarget();
		const reload = vi.fn();
		const dispose = installDynamicImportRecovery({
			target,
			storage: createStorage(),
			getHref: () => 'http://localhost:35173/app',
			now: () => 1000,
			reload
		});

		target.dispatchEvent(
			new ErrorEvent('error', {
				cancelable: true,
				message: 'Importing a module script failed.'
			})
		);

		expect(reload).toHaveBeenCalledTimes(1);
		dispose();

		const nextReload = vi.fn();
		const secondDispose = installDynamicImportRecovery({
			target,
			storage: createStorage(),
			getHref: () => 'http://localhost:35173/app',
			now: () => 1000,
			reload: nextReload
		});

		const rejection = new Event('unhandledrejection', { cancelable: true }) as PromiseRejectionEvent;
		Object.defineProperty(rejection, 'reason', {
			value: new Error('Failed to fetch dynamically imported module')
		});

		target.dispatchEvent(rejection);

		expect(nextReload).toHaveBeenCalledTimes(1);
		secondDispose();
	});
});
