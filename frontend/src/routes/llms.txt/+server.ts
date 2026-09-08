import { text } from '@sveltejs/kit';
import { about } from '$lib/data/about';
import { projects } from '$lib/data/projects';
import { SITE_URL } from '$lib/site';
import type { RequestHandler } from './$types';

// Prerendered to build/llms.txt at `npm run build` time and served by the Go
// file server. Regenerates on every deploy from the same draft-filtered data
// loaders as the site, so vault changes propagate automatically.
export const prerender = true;

interface Entry {
	slug: string;
	title: string;
	summary: string;
}

export const GET: RequestHandler = () => {
	const lines = [
		'# macuahuitl',
		'',
		'> A personal site and digital garden.',
		'',
		'Every project is available as raw markdown: append `.md` to the',
		`slash-less URL (e.g. \`${SITE_URL}/projects/example.md\`), or send \`Accept: text/markdown\``,
		'to any project URL. The home page is available at',
		`\`${SITE_URL}/about.md\`, the project index at \`${SITE_URL}/projects.md\`,`,
		'and all of these URLs also respond to `Accept: text/markdown`. Drafts are never served.'
	];
	const sections: [heading: string, prefix: string, entries: Entry[]][] = [
		['Projects', 'projects', projects]
	];
	lines.push('', '## About', '', `- [${about.title}](${SITE_URL}/about.md): ${about.summary}`);
	for (const [heading, prefix, entries] of sections) {
		if (entries.length === 0) continue;
		lines.push('', `## ${heading}`, '');
		for (const entry of entries) {
			lines.push(`- [${entry.title}](${SITE_URL}/${prefix}/${entry.slug}.md): ${entry.summary}`);
		}
	}
	return text(lines.join('\n') + '\n');
};
