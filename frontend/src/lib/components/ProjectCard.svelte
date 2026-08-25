<script lang="ts">
	import type { Project } from '$lib/data/projects';
	import StatusBadge from './StatusBadge.svelte';

	let { project }: { project: Project } = $props();
</script>

<article class="card">
	<div class="top">
		<h2><a href="/projects/{project.slug}">{project.title}</a></h2>
		<StatusBadge status={project.status} />
	</div>
	<p class="summary">{project.summary}</p>
	{#if project.tech.length}
		<ul class="tech">
			{#each project.tech as item (item)}
				<li>{item}</li>
			{/each}
		</ul>
	{/if}
</article>

<style>
	.card {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		padding: 1rem;
		border: 1px solid var(--muted);
		border-radius: 6px;
	}

	.top {
		display: flex;
		justify-content: space-between;
		align-items: baseline;
		gap: 0.5rem;
	}

	h2 {
		font-size: 1.1rem;
		line-height: 1.25;
	}

	h2 a {
		text-decoration: none;
	}

	h2 a:hover {
		text-decoration: underline;
	}

	.summary {
		font-size: 0.9em;
	}

	.tech {
		display: flex;
		flex-wrap: wrap;
		gap: 0.375rem;
		list-style: none;
		padding: 0;
		margin-top: auto;
	}

	.tech li {
		font-size: 0.75em;
		padding: 0.1em 0.6em;
		border: 1px solid var(--muted);
		border-radius: 999px;
		color: var(--muted);
	}
</style>
