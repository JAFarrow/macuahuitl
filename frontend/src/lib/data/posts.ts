import type { Component } from 'svelte';

export interface PostMetadata {
	title: string;
	created: string;
	modified?: string;
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

// YAML parses `created: YYYY-MM-DD` into a Date; normalize back to YYYY-MM-DD
function toDateString(created: string): string {
	return new Date(created).toISOString().slice(0, 10);
}

export const posts: Post[] = Object.entries(modules)
	.map(([path, mod]) => ({
			slug: slugFromPath(path),
			...mod.metadata,
			created: toDateString(mod.metadata.created),
			modified: mod.metadata.modified ? toDateString(mod.metadata.modified) : undefined
	}))
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
