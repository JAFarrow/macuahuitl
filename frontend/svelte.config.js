import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import { mdsvex } from 'mdsvex';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	extensions: ['.svelte', '.md'],
	preprocess: [vitePreprocess(), mdsvex({ extensions: ['.md'] })],
	kit: {
		adapter: adapter(),
		prerender: {
			handleHttpError: ({ status, path, referrer, referenceType }) => {
				// Raw markdown at /posts/*.md will be served by the Go server,
				// not the frontend - so 404s on those paths are expected.
				if (status === 404 && path.endsWith('.md')) return;
				throw new Error(
					`${status} ${path}${referrer ? ` (${referenceType} from ${referrer})` : ''}`
				);
			}
		}
	}
};

export default config;
