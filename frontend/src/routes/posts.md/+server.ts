import { text } from '@sveltejs/kit';
import { posts } from '$lib/data/posts';
import { SITE_URL } from '$lib/site';
import type { RequestHandler } from './$types';

// Prerendered to build/posts.md at `npm run build` time and served by the Go
// server. Mirrors the `/posts/` HTML index as markdown for agents.
export const prerender = true;

export const GET: RequestHandler = () => {
	const lines = ['# Posts', ''];
	for (const post of posts) {
		lines.push(
			`- [${post.title}](${SITE_URL}/posts/${post.slug}.md): ${post.created} — ${post.summary}`
		);
	}
	return text(lines.join('\n') + '\n');
};
