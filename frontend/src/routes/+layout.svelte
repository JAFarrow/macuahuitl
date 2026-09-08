<script lang="ts">
	import '../app.css';
	import type { Snippet } from 'svelte';
	import { about } from '$lib/data/about';

	let { children }: { children: Snippet } = $props();
</script>

<div class="site">
	<header>
		<a class="site-title" href="/">Justin</a>
	</header>

	<main>
		{@render children()}
	</main>

	{#if about.contact.length}
		<footer>
			{#each about.contact as link (link.label)}
				<a href={link.url} target="_blank" rel="noopener noreferrer">{link.label}</a>
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
		margin-bottom: 2rem;
	}

	.site-title {
		font-weight: bold;
		text-decoration: none;
		color: var(--fg);
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
