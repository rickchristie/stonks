import { apiPost } from './client';
import type { ApiResponse } from '$lib/types/api';

export type HelloResp = ApiResponse & {
	message: string;
	databaseMessage: string;
	noteCount: number;
};

export function getHello(): Promise<HelloResp> {
	return apiPost<HelloResp>('/api/hello', {});
}
