import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	const response = await resolve(event);

	const contentType = response.headers.get('content-type') || '';
	if (contentType.includes('text/html')) {
		response.headers.set('Cache-Control', 'no-store');
	}

	response.headers.set('X-Frame-Options', 'DENY');
	response.headers.set('X-Content-Type-Options', 'nosniff');
	response.headers.set('Referrer-Policy', 'strict-origin-when-cross-origin');

	const apiBase = import.meta.env.VITE_API_BASE || 'http://localhost:8080';
	response.headers.set(
		'Content-Security-Policy',
		[
			"default-src 'self'",
			"script-src 'self' 'unsafe-inline'",
			"style-src 'self' 'unsafe-inline'",
			"font-src 'self' data:",
			"img-src 'self' data: blob:",
			`connect-src 'self' ${apiBase} ${apiBase.replace(/^http/, 'ws')}`
		].join('; ')
	);

	return response;
};
