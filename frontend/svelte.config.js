import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import { mdsvex } from 'mdsvex';

// Obsidian authors image embeds as bare filenames (e.g. ![](40k.png)) which the
// vault resolves via attachmentFolderPath. Rewrite them to /attachments/*, which
// the Go server serves from content/attachments/.
function remarkAttachments() {
	return (tree) => {
		(function visit(node) {
			if (node.type === 'image' && node.url && !/^(?:[a-z]+:|\/)/i.test(node.url)) {
				node.url = `/attachments/${node.url.split('/').pop()}`;
			}
			node.children?.forEach(visit);
		})(tree);
	};
}

// Open external links (and mailto/tel) in a new tab.
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

/** @type {import('@sveltejs/kit').Config} */
const config = {
	extensions: ['.svelte', '.md'],
	preprocess: [
		vitePreprocess(),
		mdsvex({ extensions: ['.md'], remarkPlugins: [remarkAttachments, remarkExternalLinks] })
	],
	kit: {
		adapter: adapter(),
		prerender: {
			handleHttpError: ({ status, path, referrer, referenceType }) => {
				// Raw markdown at /projects/*.md and /cv/*.md, and files under
				// /attachments/, will be served by the Go server, not the frontend -
				// so 404s on those paths are expected.
				// trailingSlash: 'always' makes the prerenderer report *.md URLs with a
				// trailing slash, so tolerate both forms.
				if (
					status === 404 &&
					(path.endsWith('.md') || path.endsWith('.md/') || path.startsWith('/attachments/'))
				)
					return;
				throw new Error(
					`${status} ${path}${referrer ? ` (${referenceType} from ${referrer})` : ''}`
				);
			}
		}
	}
};

export default config;
