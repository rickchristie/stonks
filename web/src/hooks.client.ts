import type { ClientInit, HandleClientError } from '@sveltejs/kit';

import {
	installDynamicImportRecovery,
	shouldReloadForDynamicImportFailure
} from '$lib/utils/dynamic-import-recovery';

export const init: ClientInit = () => {
	installDynamicImportRecovery();
};

export const handleError: HandleClientError = ({ error, event, message }) => {
	const targetHref = event.url.href;
	if (shouldReloadForDynamicImportFailure(window.sessionStorage, targetHref, Date.now(), error)) {
		window.location.assign(targetHref);
	}

	return { message };
};
