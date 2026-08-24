<script lang="ts">
	import { postComponents } from '$lib/data/posts';
	import type { PageProps } from './$types';

	let { data }: PageProps = $props();
	const Post = $derived(postComponents[data.slug]);
</script>

<svelte:head>
	<title>{data.metadata.title} — Justin</title>
	<link rel="alternate" type="text/markdown" href="/posts/{data.slug}.md" />
</svelte:head>

<article class="prose">
	<h1>{data.metadata.title}</h1>
	<Post />
</article>
