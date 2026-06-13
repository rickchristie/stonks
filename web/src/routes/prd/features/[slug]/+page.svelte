<script lang="ts">
	let { data } = $props();
</script>

<svelte:head>
	<title>{data.doc.title} - PRD</title>
</svelte:head>

<article class="doc">
	<header>
		<p>Feature</p>
		<h1>{data.doc.title}</h1>
		<span>{data.doc.summary}</span>
	</header>

	<ul>
		{#each data.doc.body as item}
			<li>{item}</li>
		{/each}
	</ul>
</article>

<style>
	.doc,
	header,
	ul {
		display: grid;
		gap: 14px;
	}

	header p,
	h1,
	ul {
		margin: 0;
	}

	header p {
		color: rgb(var(--color-muted));
		font-family: var(--font-mono);
		font-size: 12px;
		text-transform: uppercase;
	}

	h1 {
		font-size: clamp(36px, 6vw, 70px);
		line-height: 0.98;
	}

	header span,
	li {
		color: rgb(var(--color-muted));
		font-size: 18px;
		line-height: 1.6;
	}
</style>
