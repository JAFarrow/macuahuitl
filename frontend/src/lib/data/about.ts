import type { Component } from 'svelte';

export interface AboutMetadata {
	title: string;
	summary: string;
	modified?: string;
	draft: boolean;
}

// YAML parses `modified: YYYY-MM-DD` into a Date; normalize back to YYYY-MM-DD
function toDateString(d: string): string {
	return new Date(d).toISOString().slice(0, 10);
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

export const about: AboutMetadata = {
	...mod.metadata,
	modified: mod.metadata.modified ? toDateString(mod.metadata.modified) : undefined
};
export const About: Component = mod.default;
