import { dev } from '$app/environment';
import { error } from '@sveltejs/kit';

export const prerender = false;
export const ssr = false;

export function load() {
	if (!dev) throw error(404, 'Not found');
	return {};
}
