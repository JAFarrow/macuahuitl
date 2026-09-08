<script lang="ts">
	import {
		cvComponents,
		pinnedEntries,
		workEntries,
		educationEntries,
		formatPeriod,
		type CvEntry
	} from '$lib/data/cv';
	import ContentPage from '$lib/components/ContentPage.svelte';
</script>

<svelte:head>
	<title>CV — Justin</title>
	<link rel="alternate" type="text/markdown" href="/cv.md" />
</svelte:head>

<ContentPage title="CV">
	{#each pinnedEntries as entry (entry.slug)}
		{@const Body = cvComponents[entry.slug]}
		{#if entry.kind === 'skills'}
			<h2>{entry.title}</h2>
		{/if}
		<Body />
	{/each}

	<h2>Experience</h2>
	<div class="timeline">
		{#each workEntries as entry (entry.slug)}
			{@render row(entry)}
		{/each}
	</div>

	<h2>Education</h2>
	<div class="timeline">
		{#each educationEntries as entry (entry.slug)}
			{@render row(entry)}
		{/each}
	</div>
</ContentPage>

{#snippet row(entry: CvEntry)}
	<article class="entry" class:current={!entry.end}>
		<h3><a href="/cv/{entry.slug}/">{entry.title}</a></h3>
		<p class="meta">{entry.org} · {formatPeriod(entry.start, entry.end)}</p>
		<p class="summary">{entry.summary}</p>
	</article>
{/snippet}

<style>
	.timeline {
		margin-top: 1rem;
		border-left: 2px solid var(--muted);
		padding-left: 1.5rem;
	}

	.entry {
		position: relative;
	}

	.entry + .entry {
		margin-top: 2rem;
	}

	/* Dot marker on the spine: grey hollow circle; only the current role
	   is filled with the accent colour. */
	.entry::before {
		content: '';
		position: absolute;
		left: calc(-1.5rem - 2px - 0.4375rem);
		top: 0.45rem;
		width: 0.875rem;
		height: 0.875rem;
		border-radius: 50%;
		background: var(--bg);
		border: 2px solid var(--muted);
	}

	.entry.current::before {
		background: var(--accent);
		border-color: var(--accent);
	}

	.entry h3 {
		font-size: 1.2rem;
	}

	.entry h3 a {
		color: var(--fg);
		text-decoration: none;
	}

	.entry h3 a:hover {
		color: var(--accent);
		text-decoration: underline;
	}

	.meta {
		color: var(--muted);
		font-size: 0.875em;
		margin-top: 0.25rem;
	}

	.summary {
		margin-top: 0.5rem;
	}
</style>
