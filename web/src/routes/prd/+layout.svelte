<script lang="ts">
	import PrdShell from '$lib/components/prd/PrdShell.svelte';

	let { children } = $props();
</script>

<PrdShell>
	{@render children()}
</PrdShell>
