<script lang="ts">
	import '../app.css';
	import type { Snippet } from 'svelte';
	import { page } from '$app/state';
	import { about } from '$lib/data/about';

	let { children }: { children: Snippet } = $props();
	const onProject = $derived(page.url.pathname.startsWith('/projects/'));
</script>

<div class="site">
	<header>
		<h2 class="site-title">Justin Farrow</h2>
		{#if onProject}
			<a class="backlink" href="/">← Home</a>
		{/if}
	</header>

	<main>
		{@render children()}
	</main>

	{#if about.contact.length}
		<footer>
			{#each about.contact as link (link.label)}
				<a href={link.url} target="_blank" rel="noopener">{link.label}</a>
			{/each}
		</footer>
	{/if}
</div>

<style>
	.site {
		max-width: var(--measure);
		margin: 0 auto;
		padding: 2rem 1rem;
		min-height: 100vh;
		display: flex;
		flex-direction: column;
	}

	main {
		flex: 1;
	}

	header {
		display: flex;
		justify-content: space-between;
		align-items: baseline;
		margin-bottom: 2rem;
	}

	.site-title {
		font-weight: bold;
		text-decoration: none;
		color: var(--fg);
	}

	.backlink {
		color: var(--muted);
		text-decoration: none;
	}

	.backlink:hover {
		text-decoration: underline;
	}

	footer {
		display: flex;
		justify-content: center;
		flex-wrap: wrap;
		gap: 1rem;
		margin-top: 2rem;
		padding-top: 1rem;
		border-top: 1px solid var(--muted);
		font-size: 0.875em;
	}
</style>
