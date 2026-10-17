<script lang="ts">
	import { about, About } from '$lib/data/about';
	import { projects } from '$lib/data/projects';
	import ContentPage from '$lib/components/ContentPage.svelte';
	import ProjectCard from '$lib/components/ProjectCard.svelte';
</script>

<svelte:head>
	<title>Justin Farrow</title>
	<link rel="alternate" type="text/markdown" href="/about.md" />
</svelte:head>

<ContentPage title={about.title}>
	<About />
</ContentPage>

<section>
	<h1>Projects</h1>
	<div class="card-grid">
		{#each projects as project (project.slug)}
			<ProjectCard {project} />
		{/each}
	</div>
</section>

<style>
	section {
		margin-top: 2rem;
	}

	h1 {
		margin-bottom: .5rem;
	}
</style>
