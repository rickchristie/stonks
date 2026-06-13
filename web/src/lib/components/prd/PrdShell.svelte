<script lang="ts">
	import { decisionDocs, featureDocs } from '$lib/prd/data';

	let { children } = $props();
</script>

<div class="prd-shell">
	<aside>
		<a class="brand" href="/prd">PRD</a>
		<nav aria-label="PRD navigation">
			<p>Features</p>
			{#each featureDocs as doc}
				<a href={doc.path}>{doc.title}</a>
			{/each}
			<p>Decisions</p>
			{#each decisionDocs as doc}
				<a href={doc.path}>{doc.title}</a>
			{/each}
			<p>Scenarios</p>
			<a href="/prd/storybook">Storybook</a>
		</nav>
	</aside>
	<section class="content">
		{@render children()}
	</section>
</div>

<style>
	.prd-shell {
		display: grid;
		grid-template-columns: 260px minmax(0, 1fr);
		min-height: 100vh;
		background: rgb(var(--color-bg));
	}

	aside {
		border-right: 1px solid rgb(var(--color-line));
		padding: 24px 20px;
	}

	.brand {
		display: inline-block;
		margin-bottom: 24px;
		color: rgb(var(--color-ink));
		font-size: 28px;
		font-weight: 650;
		text-decoration: none;
	}

	nav {
		display: grid;
		gap: 10px;
	}

	nav p {
		margin: 18px 0 2px;
		color: rgb(var(--color-muted));
		font-family: var(--font-mono);
		font-size: 11px;
		text-transform: uppercase;
	}

	nav a {
		color: rgb(var(--color-muted));
		text-decoration: none;
		line-height: 1.35;
	}

	nav a:hover {
		color: rgb(var(--color-accent));
	}

	.content {
		width: min(920px, calc(100vw - 320px));
		padding: 42px;
	}

	@media (max-width: 820px) {
		.prd-shell {
			grid-template-columns: 1fr;
		}

		aside {
			border-right: 0;
			border-bottom: 1px solid rgb(var(--color-line));
		}

		.content {
			width: auto;
			padding: 28px 20px;
		}
	}
</style>
