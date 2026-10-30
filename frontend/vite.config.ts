import adapter from '@sveltejs/adapter-static';
import { sveltekit } from '@sveltejs/kit/vite';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import { mdsvex } from 'mdsvex';
import { defineConfig } from 'vite';

function remarkExternalLinks() {
	return (tree: any) => {
		(function visit(node: any) {
			if (node.type === 'link' && node.url && /^(?:https?:|mailto:|tel:|\/\/)/i.test(node.url)) {
				node.data ||= {};
				node.data.hProperties ||= {};
				node.data.hProperties.target = '_blank';
				node.data.hProperties.rel = 'noopener noreferrer';
			}
			node.children?.forEach(visit);
		})(tree);
	};
}

export default defineConfig({
	plugins: [
		sveltekit({
			extensions: ['.svelte', '.md'],
			preprocess: [
				vitePreprocess(),
				mdsvex({ extensions: ['.md'], remarkPlugins: [remarkExternalLinks] })
			],
			adapter: adapter(),
			prerender: {
				handleHttpError: ({ status, path, referrer, referenceType }) => {
					if (status === 404 && (path.endsWith('.md') || path.endsWith('.md/'))) return;
					throw new Error(
						`${status} ${path}${referrer ? ` (${referenceType} from ${referrer})` : ''}`
					);
				}
			}
		})
	],
	server: {
		fs: {
			allow: ['..']
		},
		proxy: {
			'/api': 'http://localhost:8080'
		}
	}
});
