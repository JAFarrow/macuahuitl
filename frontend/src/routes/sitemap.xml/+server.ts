import { about } from '$lib/data/about';
import { cvEntries } from '$lib/data/cv';
import { projects } from '$lib/data/projects';
import { SITE_URL } from '$lib/site';
import type { RequestHandler } from './$types';

// Prerendered to build/sitemap.xml at `npm run build` time and served by the
// Go file server, same as llms.txt. Regenerates on every deploy from the same
// draft-filtered data loaders as the site, so drafts never appear here either.
export const prerender = true;

// Dates are normalized to YYYY-MM-DD by the loaders, so lexicographic max
// is a correct date max.
const latest = (items: { created: string; modified?: string }[]): string | undefined =>
	items.reduce<string | undefined>((max, i) => {
		const d = i.modified ?? i.created;
		return max === undefined || d > max ? d : max;
	}, undefined);

// Sitemap URLs use the canonical trailing-slash form; the slash-less URLs
// 301 to these.
export const GET: RequestHandler = () => {
	const urls: { loc: string; lastmod?: string }[] = [
		{ loc: '/', lastmod: about.modified },
		{ loc: '/projects/', lastmod: latest(projects) },
		{ loc: '/cv/', lastmod: latest(cvEntries) },
		...projects.map((project) => ({ loc: `/projects/${project.slug}/`, lastmod: project.modified ?? project.created })),
		...cvEntries.map((entry) => ({ loc: `/cv/${entry.slug}/`, lastmod: entry.modified ?? entry.created }))
	];
	const entries = urls
		.map(
			({ loc, lastmod }) =>
				`\t<url>\n\t\t<loc>${SITE_URL}${loc}</loc>\n` +
				(lastmod ? `\t\t<lastmod>${lastmod}</lastmod>\n` : '') +
				'\t</url>'
		)
		.join('\n');
	const body =
		'<?xml version="1.0" encoding="UTF-8"?>\n' +
		'<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">\n' +
		entries +
		'\n</urlset>\n';
	return new Response(body, {
		headers: { 'Content-Type': 'application/xml; charset=utf-8' }
	});
};
