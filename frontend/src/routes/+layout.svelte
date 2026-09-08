<script lang="ts">
	import '../app.css';
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<div class="site">
	<header>
		<a class="site-title" href="/">Justin</a>
		<nav>
			<a href="/">Home</a>
			<a href="/projects">Projects</a>
			<a href="/cv">CV</a>
		</nav>
	</header>

	<main>
		{@render children()}
	</main>
</div>

<style>
	.site {
		max-width: var(--measure);
		margin: 0 auto;
		padding: 2rem 1rem;
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

	nav a {
		margin-left: 1rem;
	}
</style>
