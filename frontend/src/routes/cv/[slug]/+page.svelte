<script lang="ts">
	import { cvComponents, formatPeriod } from '$lib/data/cv';
	import ContentPage from '$lib/components/ContentPage.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const Entry = $derived(cvComponents[data.slug]);
	const period = $derived(formatPeriod(data.metadata.start, data.metadata.end));
</script>

<svelte:head>
	<title>{data.metadata.title} — Justin</title>
	<link rel="alternate" type="text/markdown" href="/cv/{data.slug}.md" />
</svelte:head>

<ContentPage title={data.metadata.title}>
	{#snippet meta()}
		{#if data.metadata.org}
			<span>{data.metadata.org}</span>
		{/if}
		{#if period}
			<span>{period}</span>
		{/if}
	{/snippet}
	<Entry />
</ContentPage>
