import { defineConfig } from 'vitest/config';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
	plugins: [svelte({ hot: false })],
	resolve: {
		conditions: ['browser']
	},
	test: {
		include: ['src/**/*.test.ts'],
		environment: 'jsdom',
		globals: true,
		alias: {
			$lib: '/src/lib',
			'$app/environment': '/src/test/app-environment.ts',
			'$app/navigation': '/src/test/app-navigation.ts'
		}
	}
});
