import type { ApiResponse } from '$lib/types/api';

const API_BASE = import.meta.env.VITE_API_BASE ?? 'http://localhost:8080';

export interface ApiResult<T> {
	data: T | null;
	networkError: boolean;
	status?: number;
	error: string | null;
}

export class ApiTransportError extends Error {
	readonly endpoint: string;
	readonly status?: number;
	readonly responseBody?: string;

	constructor(endpoint: string, message: string, status?: number, responseBody?: string) {
		super(message);
		this.name = 'ApiTransportError';
		this.endpoint = endpoint;
		this.status = status;
		this.responseBody = responseBody;
	}
}

export async function apiPost<T extends ApiResponse>(
	endpoint: string,
	body: object,
	options: { returnError: true }
): Promise<ApiResult<T>>;
export async function apiPost<T extends ApiResponse>(
	endpoint: string,
	body?: object,
	options?: { returnError?: false }
): Promise<T>;
export async function apiPost<T extends ApiResponse>(
	endpoint: string,
	body: object = {},
	options: { returnError?: boolean } = {}
): Promise<T | ApiResult<T>> {
	const returnError = options.returnError ?? false;

	try {
		const resp = await fetch(`${API_BASE}${endpoint}`, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			credentials: 'include',
			body: JSON.stringify(body)
		});

		if (!resp.ok) {
			const responseBody = await resp.text().catch(() => '');
			const message = `HTTP ${resp.status} from ${endpoint}`;
			if (returnError) {
				return {
					data: null,
					networkError: false,
					status: resp.status,
					error: responseBody || message
				};
			}
			throw new ApiTransportError(endpoint, message, resp.status, responseBody);
		}

		const data = (await resp.json()) as T;
		if (returnError) {
			return { data, networkError: false, status: resp.status, error: null };
		}
		return data;
	} catch (err) {
		if (err instanceof ApiTransportError) throw err;
		const message = err instanceof Error ? err.message : String(err);
		if (returnError) {
			return { data: null, networkError: true, error: message };
		}
		throw new ApiTransportError(endpoint, `Network error from ${endpoint}: ${message}`);
	}
}

export async function apiPostPublic<T extends ApiResponse>(
	endpoint: string,
	body: object = {}
): Promise<{ data: T | null; error: string | null }> {
	const result = await apiPost<T>(endpoint, body, { returnError: true });
	if (result.error) return { data: null, error: result.error };
	return { data: result.data, error: null };
}

export async function checkHealth(): Promise<boolean> {
	try {
		const controller = new AbortController();
		const timeoutId = setTimeout(() => controller.abort(), 3000);
		const resp = await fetch(`${API_BASE}/health-check`, {
			method: 'GET',
			signal: controller.signal
		});
		clearTimeout(timeoutId);
		return resp.ok;
	} catch {
		return false;
	}
}
