<script lang="ts">
	import { projectComponents } from '$lib/data/projects';
	import ContentPage from '$lib/components/ContentPage.svelte';
	import StatusBadge from '$lib/components/StatusBadge.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const Project = $derived(projectComponents[data.slug]);
</script>

<svelte:head>
	<title>{data.metadata.title} — Justin</title>
	<link rel="alternate" type="text/markdown" href="/projects/{data.slug}.md" />
</svelte:head>

<a class="backlink" href="/">← Projects</a>

<ContentPage title={data.metadata.title}>
	{#snippet meta()}
		<StatusBadge status={data.metadata.status} />
		<time>{data.metadata.created}</time>
		{#if data.metadata.repo}
			<a href={data.metadata.repo}>Repo</a>
		{/if}
		{#if data.metadata.link}
			<a href={data.metadata.link}>Link</a>
		{/if}
	{/snippet}
	<Project />
</ContentPage>

<style>
	.backlink {
		display: inline-block;
		margin-bottom: 1.5rem;
		color: var(--muted);
		text-decoration: none;
	}

	.backlink:hover {
		text-decoration: underline;
	}
</style>
