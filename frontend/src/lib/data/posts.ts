import type { Component } from 'svelte';

export interface PostMetadata {
	title: string;
	created: string;
	tags: string[];
	summary: string;
	status: string;
	draft: boolean;
}

export interface Post extends PostMetadata {
	slug: string;
}

interface MarkdownModule {
	metadata: PostMetadata;
	default: Component;
}

const modules = import.meta.glob<MarkdownModule>('../../../../content/posts/*.md', {
	eager: true
});

function slugFromPath(path: string): string {
	return path.split('/').pop()!.replace(/\.md$/, '');
}

export const posts: Post[] = Object.entries(modules)
	.map(([path, mod]) => ({ slug: slugFromPath(path), ...mod.metadata }))
	.filter((post) => !post.draft)
	.sort((a, b) => b.created.localeCompare(a.created));

export const postComponents: Record<string, Component> = Object.fromEntries(
	Object.entries(modules)
		.filter(([, mod]) => !mod.metadata.draft)
		.map(([path, mod]) => [slugFromPath(path), mod.default])
);

export function getPost(slug: string): Post | undefined {
	return posts.find((post) => post.slug === slug);
}
