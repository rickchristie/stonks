/**
 * Client-route chunks can fail to load after a deploy or while Vite regenerates
 * its client manifest during development. Reload the current route once for
 * this narrow failure class so the browser fetches the current asset graph
 * instead of leaving the user on SvelteKit's generic error page.
 */
const DYNAMIC_IMPORT_RECOVERY_KEY = 'stonks-dynamic-import-recovery';
const RECOVERY_WINDOW_MS = 15000;
const DYNAMIC_IMPORT_FAILURE_RE =
	/(Failed to fetch dynamically imported module|Importing a module script failed|error loading dynamically imported module)/i;

export interface StorageLike {
	getItem(key: string): string | null;
	setItem(key: string, value: string): void;
}

interface RecoveryRecord {
	href: string;
	ts: number;
}

interface RecoveryEvent {
	preventDefault?: () => void;
}

export interface DynamicImportRecoveryOptions {
	target?: EventTarget;
	storage?: StorageLike;
	getHref?: () => string;
	now?: () => number;
	reload?: () => void;
}

export function extractDynamicImportFailureMessage(reason: unknown): string {
	if (typeof reason === 'string') return reason;
	if (reason instanceof Error) return reason.message;
	if (typeof reason === 'object' && reason !== null) {
		const message = (reason as { message?: unknown }).message;
		if (typeof message === 'string') return message;
	}
	return '';
}

export function isDynamicImportFailure(reason: unknown): boolean {
	return DYNAMIC_IMPORT_FAILURE_RE.test(extractDynamicImportFailureMessage(reason));
}

export function shouldReloadForDynamicImportFailure(
	storage: StorageLike,
	href: string,
	now: number,
	reason: unknown
): boolean {
	if (!isDynamicImportFailure(reason)) return false;

	try {
		const raw = storage.getItem(DYNAMIC_IMPORT_RECOVERY_KEY);
		if (raw) {
			const previous = JSON.parse(raw) as Partial<RecoveryRecord>;
			if (
				previous.href === href &&
				typeof previous.ts === 'number' &&
				now - previous.ts < RECOVERY_WINDOW_MS
			) {
				return false;
			}
		}
	} catch {
		// Malformed session state should not block the one-shot recovery attempt.
	}

	try {
		storage.setItem(DYNAMIC_IMPORT_RECOVERY_KEY, JSON.stringify({ href, ts: now } satisfies RecoveryRecord));
	} catch {
		// Best effort only; private browsing or disabled storage should still reload.
	}

	return true;
}

export function installDynamicImportRecovery({
	target = window,
	storage = window.sessionStorage,
	getHref = () => window.location.href,
	now = () => Date.now(),
	reload = () => window.location.reload()
}: DynamicImportRecoveryOptions = {}): () => void {
	let reloadQueued = false;

	const maybeReload = (reason: unknown, event?: RecoveryEvent) => {
		if (reloadQueued) return;
		if (!shouldReloadForDynamicImportFailure(storage, getHref(), now(), reason)) return;
		reloadQueued = true;
		event?.preventDefault?.();
		reload();
	};

	const handlePreloadError = (event: Event) => {
		const preloadEvent = event as CustomEvent<unknown> & { payload?: unknown };
		maybeReload(preloadEvent.payload ?? preloadEvent.detail, preloadEvent);
	};

	const handleError = (event: Event) => {
		const errorEvent = event as ErrorEvent;
		maybeReload(errorEvent.error ?? errorEvent.message, errorEvent);
	};

	const handleUnhandledRejection = (event: Event) => {
		const rejectionEvent = event as PromiseRejectionEvent;
		maybeReload(rejectionEvent.reason, rejectionEvent);
	};

	target.addEventListener('vite:preloadError', handlePreloadError as EventListener);
	target.addEventListener('error', handleError as EventListener);
	target.addEventListener('unhandledrejection', handleUnhandledRejection as EventListener);

	return () => {
		target.removeEventListener('vite:preloadError', handlePreloadError as EventListener);
		target.removeEventListener('error', handleError as EventListener);
		target.removeEventListener('unhandledrejection', handleUnhandledRejection as EventListener);
	};
}
