<script lang="ts">
	import { projects } from '$lib/data/projects';
	import PostList from '$lib/components/PostList.svelte';
	import ProjectCard from '$lib/components/ProjectCard.svelte';
</script>

<svelte:head>
	<title>Justin</title>
</svelte:head>

<h1>Justin</h1>
<p>A personal site and digital garden. Notes here are perpetually unfinished.</p>

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
