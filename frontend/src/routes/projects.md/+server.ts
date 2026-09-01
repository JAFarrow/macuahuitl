import { text } from '@sveltejs/kit';
import { projects } from '$lib/data/projects';
import { SITE_URL } from '$lib/site';
import type { RequestHandler } from './$types';

// Prerendered to build/projects.md at `npm run build` time and served by the Go
// server. Mirrors the `/projects/` HTML index as markdown for agents.
export const prerender = true;

export const GET: RequestHandler = () => {
	const lines = ['# Projects', ''];
	for (const project of projects) {
		lines.push(
			`- [${project.title}](${SITE_URL}/projects/${project.slug}.md): ${project.status} — ${project.summary}`
		);
	}
	return text(lines.join('\n') + '\n');
};
