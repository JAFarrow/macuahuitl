<script lang="ts">
	import { about, About } from '$lib/data/about';
	import { projects } from '$lib/data/projects';
	import ContentPage from '$lib/components/ContentPage.svelte';
	import PostList from '$lib/components/PostList.svelte';
	import ProjectCard from '$lib/components/ProjectCard.svelte';
</script>

<svelte:head>
	<title>Justin</title>
	<link rel="alternate" type="text/markdown" href="/about.md" />
</svelte:head>

<ContentPage title={about.title}>
	<About />
</ContentPage>

<section>
	<h2>Posts</h2>
	<PostList />
</section>

<section>
	<h2>Projects</h2>
	<div class="card-grid">
		{#each projects as project (project.slug)}
			<ProjectCard {project} />
		{/each}
	</div>
</section>

<style>
	section {
		margin-top: 2.5rem;
	}

	h2 {
		margin-bottom: 1rem;
	}
</style>
