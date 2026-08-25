<script lang="ts">
	import type { ProjectStatus } from '$lib/data/projects';

	let { status }: { status: ProjectStatus } = $props();

	const glyphs: Record<ProjectStatus, string> = {
		seedling: '🌱',
		budding: '🌿',
		evergreen: '🌳',
		archived: '📦'
	};
</script>

<span class="badge" data-status={status}>{glyphs[status]} {status}</span>

<style>
	.badge {
		font-size: 0.75em;
		white-space: nowrap;
		padding: 0.15em 0.6em;
		border: 1px solid var(--muted);
		border-radius: 999px;
		color: var(--muted);
	}

	.badge[data-status='evergreen'] {
		color: var(--accent);
		border-color: var(--accent);
	}
</style>
