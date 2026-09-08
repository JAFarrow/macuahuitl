import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		// content/ lives at the repo root, outside frontend/
		fs: {
			allow: ['..']
		},
		// The Go server doesn't exist yet — this will 502 in dev. Expected.
		proxy: {
			'/api': 'http://localhost:8080'
		}
	}
});
