import type { Component } from 'svelte';

export interface ContactLink {
	label: string;
	url: string;
}

const CONTACT_KEYS = ['email', 'linkedin', 'github'] as const;
type ContactKey = (typeof CONTACT_KEYS)[number];

interface AboutFrontmatter extends Partial<Record<ContactKey, string>> {
	title: string;
	summary: string;
	modified?: string;
	draft: boolean;
}

export interface AboutMetadata {
	title: string;
	summary: string;
	modified?: string;
	draft: boolean;
	contact: ContactLink[];
}

function toDateString(d: string): string {
	return new Date(d).toISOString().slice(0, 10);
}

interface MarkdownModule {
	metadata: AboutFrontmatter;
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

const contact: ContactLink[] = CONTACT_KEYS.flatMap((label) => {
	const url = mod.metadata[label];
	return url ? [{ label, url }] : [];
});

export const about: AboutMetadata = {
	title: mod.metadata.title,
	summary: mod.metadata.summary,
	modified: mod.metadata.modified ? toDateString(mod.metadata.modified) : undefined,
	draft: mod.metadata.draft,
	contact
};
export const About: Component = mod.default;
