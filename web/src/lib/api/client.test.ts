import { afterEach, describe, expect, it, vi } from 'vitest';
import { ApiTransportError, apiPost, apiPostPublic, checkHealth } from './client';

afterEach(() => {
	vi.unstubAllGlobals();
	vi.useRealTimers();
});

describe('apiPost', () => {
	it('posts JSON and returns parsed body', async () => {
		const fetchMock = vi.fn().mockResolvedValue({
			ok: true,
			status: 200,
			json: async () => ({ error: '', value: 1 })
		});
		vi.stubGlobal('fetch', fetchMock);

		const got = await apiPost<{ error: string; value: number }>('/api/example', { id: 1 });

		expect(got).toEqual({ error: '', value: 1 });
		expect(fetchMock).toHaveBeenCalledWith('http://localhost:8080/api/example', {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify({ id: 1 })
		});
	});

	it('throws for non-200 transport responses', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn().mockResolvedValue({ ok: false, status: 500, text: async () => 'server exploded' })
		);

		await expect(apiPost('/api/example', {})).rejects.toThrow(ApiTransportError);
		await expect(apiPost('/api/example', {})).rejects.toMatchObject({
			endpoint: '/api/example',
			status: 500,
			responseBody: 'server exploded'
		});
	});

	it('returns transport details when requested', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn().mockResolvedValue({ ok: false, status: 503, text: async () => 'temporarily unavailable' })
		);

		const got = await apiPost('/api/example', {}, { returnError: true });

		expect(got).toEqual({
			data: null,
			networkError: false,
			status: 503,
			error: 'temporarily unavailable'
		});
	});

	it('keeps application-level errors in response data', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn().mockResolvedValue({
				ok: true,
				status: 200,
				json: async () => ({ error: 'Validation', value: 0 })
			})
		);

		const got = await apiPost<{ error: string; value: number }>('/api/example', {});

		expect(got).toEqual({ error: 'Validation', value: 0 });
	});

	it('wraps public endpoint transport failures without throwing', async () => {
		vi.stubGlobal(
			'fetch',
			vi.fn().mockResolvedValue({ ok: false, status: 429, text: async () => 'too many requests' })
		);

		await expect(apiPostPublic('/api/login', {})).resolves.toEqual({
			data: null,
			error: 'too many requests'
		});
	});

	it('checks backend health with a timeout', async () => {
		vi.useFakeTimers();
		vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true }));

		await expect(checkHealth()).resolves.toBe(true);
	});
});
