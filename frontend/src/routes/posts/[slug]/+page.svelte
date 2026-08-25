<script lang="ts">
	import { postComponents } from '$lib/data/posts';
	import ContentPage from '$lib/components/ContentPage.svelte';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const Post = $derived(postComponents[data.slug]);
</script>

<svelte:head>
	<title>{data.metadata.title} — Justin</title>
	<link rel="alternate" type="text/markdown" href="/posts/{data.slug}.md" />
</svelte:head>

<ContentPage title={data.metadata.title}>
	{#snippet meta()}
		<time>{data.metadata.created}</time>
	{/snippet}
	<Post />
</ContentPage>
