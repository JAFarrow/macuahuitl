import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import { mdsvex } from 'mdsvex';

function remarkExternalLinks() {
	return (tree) => {
		(function visit(node) {
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

const config = {
	extensions: ['.svelte', '.md'],
	preprocess: [
		vitePreprocess(),
		mdsvex({ extensions: ['.md'], remarkPlugins: [remarkExternalLinks] })
	],
	kit: {
		adapter: adapter(),
		prerender: {
			handleHttpError: ({ status, path, referrer, referenceType }) => {
				if (status === 404 && (path.endsWith('.md') || path.endsWith('.md/'))) return;
				throw new Error(
					`${status} ${path}${referrer ? ` (${referenceType} from ${referrer})` : ''}`
				);
			}
		}
	}
};

export default config;
