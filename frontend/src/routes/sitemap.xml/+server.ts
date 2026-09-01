import { about } from '$lib/data/about';
import { posts } from '$lib/data/posts';
import { projects } from '$lib/data/projects';
import type { RequestHandler } from './$types';

// Prerendered to build/sitemap.xml at `npm run build` time and served by the
// Go file server, same as llms.txt. Regenerates on every deploy from the same
// draft-filtered data loaders as the site, so drafts never appear here either.
export const prerender = true;

// Canonical origin of the deployed site; sitemap URLs must be absolute.
const SITE_URL = 'https://www.justin-farrow-dev.com';

// Sitemap URLs use the canonical trailing-slash form; the slash-less URLs
// 301 to these.
export const GET: RequestHandler = () => {
	const urls: { loc: string; lastmod?: string }[] = [
		{ loc: '/', lastmod: about.modified },
		{ loc: '/posts/' },
		{ loc: '/projects/' },
		...posts.map((post) => ({ loc: `/posts/${post.slug}/`, lastmod: post.modified ?? post.created })),
		...projects.map((project) => ({ loc: `/projects/${project.slug}/`, lastmod: project.modified ?? project.created }))
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
