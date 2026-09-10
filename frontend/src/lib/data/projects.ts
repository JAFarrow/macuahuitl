import type { Component } from 'svelte';

export type ProjectStatus = 'seedling' | 'fruited' | 'evergreen' | 'archived';

export interface ProjectMetadata {
	title: string;
	created: string;
	modified?: string;
	status: ProjectStatus;
	tech: string[];
	summary: string;
	repo?: string;
	link?: string;
	draft: boolean;
}

export interface Project extends ProjectMetadata {
	slug: string;
}

interface MarkdownModule {
	metadata: ProjectMetadata;
	default: Component;
}

const modules = import.meta.glob<MarkdownModule>('../../../../content/projects/*.md', {
	eager: true
});

function slugFromPath(path: string): string {
	return path.split('/').pop()!.replace(/\.md$/, '');
}

function toDateString(created: string): string {
	return new Date(created).toISOString().slice(0, 10);
}

export const projects: Project[] = Object.entries(modules)
	.map(([path, mod]) => ({
			slug: slugFromPath(path),
			...mod.metadata,
			created: toDateString(mod.metadata.created),
			modified: mod.metadata.modified ? toDateString(mod.metadata.modified) : undefined
	}))
	.filter((project) => !project.draft)
	.sort((a, b) => b.created.localeCompare(a.created));

export const projectComponents: Record<string, Component> = Object.fromEntries(
	Object.entries(modules)
		.filter(([, mod]) => !mod.metadata.draft)
		.map(([path, mod]) => [slugFromPath(path), mod.default])
);

export function getProject(slug: string): Project | undefined {
	return projects.find((project) => project.slug === slug);
}
