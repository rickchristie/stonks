<script lang="ts">
	import { prdIndexGroups } from '$lib/prd/data';
</script>

<svelte:head>
	<title>PRD - Stonks</title>
</svelte:head>

<article class="doc">
	<header>
		<p>Product Requirements</p>
		<h1>Stonks</h1>
		<span>Preparation notes for Rick's stock research, portfolio management, and automated trading app.</span>
	</header>

	{#each prdIndexGroups as group}
		<section>
			<h2>{group.title}</h2>
			<ul>
				{#each group.items as item}
					<li><a href={item.href}>{item.title}</a> - {item.summary}</li>
				{/each}
			</ul>
		</section>
	{/each}
</article>

<style>
	.doc {
		display: grid;
		gap: 34px;
	}

	header,
	section {
		display: grid;
		gap: 12px;
	}

	header p {
		margin: 0;
		color: rgb(var(--color-muted));
		font-family: var(--font-mono);
		font-size: 12px;
		text-transform: uppercase;
	}

	h1,
	h2,
	ul {
		margin: 0;
	}

	h1 {
		font-size: clamp(42px, 7vw, 82px);
		line-height: 0.95;
	}

	h2 {
		font-size: 28px;
	}

	header span,
	li {
		color: rgb(var(--color-muted));
		font-size: 18px;
		line-height: 1.55;
	}

	ul {
		display: grid;
		gap: 10px;
		padding-left: 0;
		list-style: none;
	}

	a {
		color: rgb(var(--color-accent));
	}
</style>
