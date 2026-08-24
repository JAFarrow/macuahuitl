<script lang="ts">
	import { posts } from '$lib/data/posts';
</script>

<ul>
	{#each posts as post (post.slug)}
		<li>
			<a href="/posts/{post.slug}">{post.title}</a>
			<span class="date">{post.created}</span>
		</li>
	{/each}
</ul>

<style>
	ul {
		list-style: none;
		padding: 0;
	}

	li + li {
		margin-top: 0.5em;
	}

	.date {
		color: var(--muted);
		font-size: 0.875em;
		margin-left: 0.5em;
	}
</style>
