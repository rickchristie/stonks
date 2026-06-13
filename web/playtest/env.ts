import { existsSync, readFileSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const playtestDir = dirname(fileURLToPath(import.meta.url));
const repoRoot = resolve(playtestDir, '../..');

loadDevEnv();

requiredEnv('DEV_BACKEND_HOST');
requiredEnv('DEV_BACKEND_PORT');
requiredEnv('DEV_WEB_HOST');
requiredEnv('DEV_WEB_PORT');
requiredEnv('VITE_API_BASE');
requiredEnv('PLAYWRIGHT_BASE_URL');
requiredEnv('PLAYWRIGHT_API_BASE');

export const baseURL = process.env.PLAYWRIGHT_BASE_URL!;
export const apiBase = process.env.PLAYWRIGHT_API_BASE!;
export const wsBase = apiBase.replace(/^http/, 'ws');

function loadDevEnv(): void {
	const devEnvPath = resolve(repoRoot, '.dev');
	if (!existsSync(devEnvPath)) {
		throw new Error(`Missing required .dev file at ${devEnvPath}`);
	}

	const lines = readFileSync(devEnvPath, 'utf8').split(/\r?\n/);
	for (const line of lines) {
		const trimmed = line.trim();
		if (!trimmed || trimmed.startsWith('#')) {
			continue;
		}

		const match = trimmed.match(/^(?:export\s+)?([A-Za-z_][A-Za-z0-9_]*)=(.*)$/);
		if (!match) {
			continue;
		}

		process.env[match[1]] = parseEnvValue(match[2]);
	}
}

function parseEnvValue(value: string): string {
	const trimmed = value.trim();
	if (trimmed.length < 2) {
		return trimmed;
	}

	const quote = trimmed[0];
	if ((quote === '"' || quote === "'") && trimmed[trimmed.length - 1] === quote) {
		return trimmed.slice(1, -1);
	}

	return trimmed;
}

function requiredEnv(key: string): void {
	if (!process.env[key]) {
		throw new Error(`Invalid .dev: ${key} is required`);
	}
}
