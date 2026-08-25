<script lang="ts">
	import type { Snippet } from 'svelte';
	import Prose from './Prose.svelte';

	let { title, meta, children }: { title: string; meta?: Snippet; children: Snippet } = $props();
</script>

<header>
	<h1>{title}</h1>
	{#if meta}
		<div class="meta">{@render meta()}</div>
	{/if}
</header>

<Prose>
	{@render children()}
</Prose>

<style>
	header {
		margin-bottom: 2rem;
	}

	h1 {
		line-height: 1.25;
	}

	.meta {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		gap: 0.5em;
		margin-top: 0.5rem;
		color: var(--muted);
		font-size: 0.875em;
	}
</style>
