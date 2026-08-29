import type { Component } from 'svelte';

export interface AboutMetadata {
	title: string;
	summary: string;
	draft: boolean;
}

interface MarkdownModule {
	metadata: AboutMetadata;
	default: Component;
}

const modules = import.meta.glob<MarkdownModule>('../../../../content/about.md', {
	eager: true
});

const mod = Object.values(modules)[0];
if (!mod) {
	throw new Error('content/about.md is missing');
}
if (mod.metadata.draft) {
	throw new Error('content/about.md is drafted');
}

export const about: AboutMetadata = mod.metadata;
export const About: Component = mod.default;
