import { text } from '@sveltejs/kit';
import { pinnedEntries, timelineEntries, formatPeriod } from '$lib/data/cv';
import type { RequestHandler } from './$types';

// Prerendered to build/cv.md at `npm run build` time and served by the Go
// server. Unlike the projects.md link index, this twin is the full
// CV stitched into one document: a CV is consumed whole, so agents get it in a
// single fetch. Draft-filtered by the same loader as the site.
export const prerender = true;

// Raw vault note bodies, inlined at build time. ?raw imports bypass the mdsvex
// transform (query-suffixed ids don't match its .md filter).
const rawNotes = import.meta.glob<string>('../../../../content/cv/*.md', {
	query: '?raw',
	import: 'default',
	eager: true
});

function slugFromPath(path: string): string {
	return path.split('/').pop()!.replace(/\.md$/, '');
}

// Naive frontmatter strip: drop the --- fence block. Not a YAML parse, same
// philosophy as the Go draft check.
function stripFrontmatter(raw: string): string {
	const lines = raw.split('\n');
	if (lines[0]?.trim() !== '---') return raw.trim();
	const end = lines.findIndex((l, i) => i > 0 && l.trim() === '---');
	return (end === -1 ? raw : lines.slice(end + 1).join('\n')).trim();
}

function body(slug: string): string {
	for (const [path, raw] of Object.entries(rawNotes)) {
		if (slugFromPath(path) === slug) return stripFrontmatter(raw);
	}
	throw new Error(`content/cv/${slug}.md raw body missing`);
}

export const GET: RequestHandler = () => {
	const lines = ['# Justin Farrow — CV', ''];
	for (const entry of pinnedEntries) {
		if (entry.kind === 'skills') lines.push('## Technical Skills', '');
		lines.push(body(entry.slug), '');
	}
	for (const entry of timelineEntries) {
		lines.push(
			`## ${entry.title} — ${entry.org}`,
			'',
			`*${formatPeriod(entry.start, entry.end)}*`,
			'',
			body(entry.slug),
			''
		);
	}
	return text(lines.join('\n').replace(/\n{3,}/g, '\n\n') + '\n');
};
