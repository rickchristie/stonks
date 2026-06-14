<svelte:head>
	<title>Storybook - Stonks</title>
</svelte:head>

<article class="storybook">
	<header>
		<p>Storybook</p>
		<h1>Hello World States</h1>
	</header>

	<section class="states" aria-label="Component scenarios">
		<div class="scenario">
			<span>Loaded</span>
			<h2>Hello, World!</h2>
			<p>This message came from PostgreSQL through the Go backend.</p>
		</div>
		<div class="scenario muted">
			<span>Loading</span>
			<h2>Loading...</h2>
			<p>Waiting for `/api/hello`.</p>
		</div>
		<div class="scenario error">
			<span>Error</span>
			<h2>InternalError</h2>
			<p>The backend could not read the database.</p>
		</div>
	</section>
</article>

<style>
	.storybook,
	header,
	.states,
	.scenario {
		display: grid;
		gap: 16px;
	}

	header p,
	h1,
	h2,
	.scenario p {
		margin: 0;
	}

	header p,
	.scenario span {
		color: rgb(var(--color-muted));
		font-family: var(--font-mono);
		font-size: 12px;
		text-transform: uppercase;
	}

	h1 {
		font-size: clamp(36px, 6vw, 70px);
		line-height: 0.98;
	}

	.states {
		grid-template-columns: repeat(3, minmax(0, 1fr));
	}

	.scenario {
		align-content: start;
		min-height: 220px;
		border: 1px solid rgb(var(--color-line));
		background: rgb(var(--color-panel));
		padding: 20px;
	}

	.scenario h2 {
		font-size: 30px;
		line-height: 1.05;
	}

	.scenario p {
		color: rgb(var(--color-muted));
		line-height: 1.5;
	}

	.error h2 {
		color: rgb(var(--color-warn));
	}

	.muted h2 {
		color: rgb(var(--color-muted));
	}

	@media (max-width: 900px) {
		.states {
			grid-template-columns: 1fr;
		}
	}
</style>
