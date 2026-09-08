import type { Component } from 'svelte';

export type CvKind = 'profile' | 'skills' | 'work' | 'education';

export interface CvMetadata {
	title: string;
	org?: string;
	kind: CvKind;
	start?: string;
	end?: string;
	summary: string;
	created: string;
	modified?: string;
	draft: boolean;
}

export interface CvEntry extends CvMetadata {
	slug: string;
}

interface MarkdownModule {
	metadata: CvMetadata;
	default: Component;
}

const modules = import.meta.glob<MarkdownModule>('../../../../content/cv/*.md', {
	eager: true
});

function slugFromPath(path: string): string {
	return path.split('/').pop()!.replace(/\.md$/, '');
}

// YAML parses `created: YYYY-MM-DD` into a Date; normalize back to YYYY-MM-DD
function toDateString(created: string): string {
	return new Date(created).toISOString().slice(0, 10);
}

// start/end are authored as YYYY-MM, which YAML leaves as a string; going
// through Date also tolerates full YYYY-MM-DD values (and the Date objects
// YAML produces for them), collapsing everything to YYYY-MM.
function toYearMonth(d: string): string {
	return new Date(d).toISOString().slice(0, 7);
}

// Pinned kinds lead the page in fixed order; dated kinds (work, education)
// follow, grouped by kind. Absent end means a current role and sorts first.
const KIND_ORDER: Record<CvKind, number> = { profile: 0, skills: 1, work: 2, education: 3 };
const sortKey = (e: CvEntry): string => `${e.end ?? '9999-12'} ${e.start ?? ''}`;

export const cvEntries: CvEntry[] = Object.entries(modules)
	.map(([path, mod]) => ({
		slug: slugFromPath(path),
		...mod.metadata,
		start: mod.metadata.start ? toYearMonth(mod.metadata.start) : undefined,
		end: mod.metadata.end ? toYearMonth(mod.metadata.end) : undefined,
		created: toDateString(mod.metadata.created),
		modified: mod.metadata.modified ? toDateString(mod.metadata.modified) : undefined
	}))
	.filter((entry) => !entry.draft)
	.sort((a, b) => {
		const kind = KIND_ORDER[a.kind] - KIND_ORDER[b.kind];
		return kind !== 0 ? kind : sortKey(b).localeCompare(sortKey(a));
	});

export const pinnedEntries = cvEntries.filter((e) => e.kind === 'profile' || e.kind === 'skills');
export const timelineEntries = cvEntries.filter((e) => e.kind === 'work' || e.kind === 'education');
export const workEntries = cvEntries.filter((e) => e.kind === 'work');
export const educationEntries = cvEntries.filter((e) => e.kind === 'education');

export const cvComponents: Record<string, Component> = Object.fromEntries(
	Object.entries(modules)
		.filter(([, mod]) => !mod.metadata.draft)
		.map(([path, mod]) => [slugFromPath(path), mod.default])
);

export function getCvEntry(slug: string): CvEntry | undefined {
	return cvEntries.find((entry) => entry.slug === slug);
}

const MONTHS = [
	'January',
	'February',
	'March',
	'April',
	'May',
	'June',
	'July',
	'August',
	'September',
	'October',
	'November',
	'December'
];

export function formatYearMonth(ym: string): string {
	const [y, m] = ym.split('-').map(Number);
	return `${MONTHS[m - 1]} ${y}`;
}

export function formatPeriod(start?: string, end?: string): string {
	if (!start) return '';
	return `${formatYearMonth(start)} – ${end ? formatYearMonth(end) : 'Present'}`;
}
