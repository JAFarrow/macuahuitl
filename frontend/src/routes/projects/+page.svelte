<script lang="ts">
	import { projects } from '$lib/data/projects';
	import ProjectCard from '$lib/components/ProjectCard.svelte';
</script>

<svelte:head>
	<title>Projects — Justin</title>
	<link rel="alternate" type="text/markdown" href="/projects.md" />
</svelte:head>

<h1>Projects</h1>

<div class="card-grid">
	{#each projects as project (project.slug)}
		<ProjectCard {project} />
	{/each}
</div>

<style>
	h1 {
		margin-bottom: 1.5rem;
	}
</style>
