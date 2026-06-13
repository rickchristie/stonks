<script lang="ts">
	import { dev } from '$app/environment';
	import { onMount } from 'svelte';
	import { getHello, type HelloResp } from '$lib/api/hello';

	let response = $state<HelloResp | null>(null);
	let error = $state('');
	let busy = $state(false);

	onMount(() => {
		void loadHello();
	});

	async function loadHello() {
		if (busy) return;
		busy = true;
		error = '';
		try {
			const resp = await getHello();
			if (resp.error) {
				error = resp.error;
				response = null;
				return;
			}
			response = resp;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unknown error';
			response = null;
		} finally {
			busy = false;
		}
	}
</script>

<svelte:head>
	<title>App Template</title>
</svelte:head>

<main class="app-shell">
	<header class="topbar">
		<div>
			<p class="eyebrow">Agentic App Template</p>
			<h1>Hello World</h1>
		</div>
		{#if dev}
			<nav aria-label="Template docs">
				<a href="/prd">PRD</a>
				<a href="/prd/storybook">Storybook</a>
			</nav>
		{/if}
	</header>

	<section class="workspace" aria-label="Hello world verification">
		<section class="hello-panel">
			<p class="label">Backend API</p>
			{#if response}
				<h2 data-testid="hello-message">{response.message}</h2>
			{:else if busy}
				<h2>Loading...</h2>
			{:else}
				<h2>Not loaded</h2>
			{/if}
			<p class="copy">
				This page calls <code>/api/hello</code>. The backend responds only after reading the seeded PostgreSQL row.
			</p>
			<button onclick={() => void loadHello()} disabled={busy}>
				{busy ? 'Calling...' : 'Call API'}
			</button>
			{#if error}
				<p class="error" role="alert">{error}</p>
			{/if}
		</section>

		<section class="db-panel" aria-label="Database response">
			<p class="label">Database Proof</p>
			{#if response}
				<p class="db-message" data-testid="database-message">{response.databaseMessage}</p>
				<dl>
					<div>
						<dt>Active rows</dt>
						<dd data-testid="note-count">{response.noteCount}</dd>
					</div>
					<div>
						<dt>Local config</dt>
						<dd>.dev required</dd>
					</div>
					<div>
						<dt>Docs</dt>
						<dd>/prd dev-only</dd>
					</div>
				</dl>
			{:else}
				<p class="db-message muted">Waiting for the API response.</p>
			{/if}
		</section>
	</section>
</main>

<style>
	.app-shell {
		display: grid;
		gap: 30px;
		width: min(1120px, calc(100vw - 32px));
		margin: 0 auto;
		padding: 30px 0 54px;
	}

	.topbar {
		display: flex;
		align-items: end;
		justify-content: space-between;
		gap: 24px;
		border-bottom: 1px solid rgb(var(--color-line));
		padding-bottom: 22px;
	}

	.eyebrow,
	.label,
	dt,
	nav {
		font-family: var(--font-mono);
	}

	.eyebrow,
	.label {
		margin: 0;
		color: rgb(var(--color-muted));
		font-size: 12px;
		text-transform: uppercase;
	}

	h1,
	h2,
	p {
		margin: 0;
	}

	h1 {
		font-size: clamp(42px, 7vw, 96px);
		font-weight: 650;
		line-height: 0.92;
	}

	nav {
		display: flex;
		gap: 16px;
		color: rgb(var(--color-muted));
		font-size: 13px;
	}

	nav a {
		text-decoration: none;
	}

	.workspace {
		display: grid;
		grid-template-columns: minmax(0, 1.08fr) minmax(300px, 0.92fr);
		gap: 28px;
		align-items: stretch;
	}

	.hello-panel,
	.db-panel {
		display: grid;
		align-content: start;
		gap: 20px;
		background: rgb(var(--color-panel));
		border: 1px solid rgb(var(--color-line));
		padding: clamp(22px, 4vw, 42px);
		min-height: 390px;
	}

	.hello-panel {
		border-left: 8px solid rgb(var(--color-accent));
	}

	h2 {
		font-size: clamp(38px, 6vw, 76px);
		font-weight: 650;
		line-height: 0.98;
	}

	.copy,
	.db-message {
		max-width: 640px;
		color: rgb(var(--color-muted));
		font-size: 19px;
		line-height: 1.52;
	}

	code {
		font-family: var(--font-mono);
		font-size: 0.88em;
		color: rgb(var(--color-accent));
	}

	button {
		justify-self: start;
		border: 1px solid rgb(var(--color-accent));
		background: rgb(var(--color-accent));
		color: white;
		padding: 11px 16px;
		cursor: pointer;
	}

	button:disabled {
		cursor: not-allowed;
		opacity: 0.55;
	}

	.error {
		color: rgb(var(--color-warn));
		font-family: var(--font-mono);
		font-size: 13px;
	}

	.muted {
		color: rgb(var(--color-muted));
	}

	dl {
		display: grid;
		gap: 1px;
		margin: 8px 0 0;
		background: rgb(var(--color-line));
	}

	dl div {
		display: grid;
		grid-template-columns: minmax(110px, 0.44fr) minmax(0, 1fr);
		gap: 16px;
		background: rgb(var(--color-panel));
		padding: 14px 0;
	}

	dt {
		color: rgb(var(--color-muted));
		font-size: 12px;
		text-transform: uppercase;
	}

	dd {
		margin: 0;
		color: rgb(var(--color-ink));
	}

	@media (max-width: 820px) {
		.topbar,
		.workspace {
			display: grid;
		}

		.workspace {
			grid-template-columns: 1fr;
		}
	}
</style>
