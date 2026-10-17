<script lang="ts">
	import type { Snippet } from 'svelte';

	let { children }: { children: Snippet } = $props();
</script>

<svelte:head>
	<title>For Miss Jordaan</title>
	<meta name="robots" content="noindex, nofollow" />
	<meta name="theme-color" content="#f3d9b1" />
</svelte:head>

<div class="beloved">
	<div class="clouds" aria-hidden="true"></div>
	<div class="grain" aria-hidden="true"></div>
	<div class="content">
		{@render children()}
	</div>
</div>

<style>
	:global(body) {
		margin: 0;
		background: #f3d9b1;
	}

	.beloved {
		--cream: #fbf1dd;
		--sand: #f3d9b1;
		--clay: #c4683a;
		--ink: #3b2416;
		--dusk: #6b3f5e;
		--sun: #e8923c;
		--serif: Georgia, 'Iowan Old Style', 'Times New Roman', serif;
		position: relative;
		min-height: 100vh;
		min-height: 100dvh;
		color: var(--ink);
		font-family: var(--serif);
		line-height: 1.55;
		background: linear-gradient(
			180deg,
			#6b3f5e 0%,
			#c4683a 38%,
			#e8923c 62%,
			#f3d9b1 100%
		);
		background-attachment: fixed;
		overflow-x: hidden;
	}

	.clouds {
		position: fixed;
		inset: 0;
		pointer-events: none;
		background:
			radial-gradient(ellipse 28% 7% at 20% 22%, rgb(255 244 225 / 0.55), transparent 70%),
			radial-gradient(ellipse 34% 8% at 78% 30%, rgb(255 244 225 / 0.45), transparent 70%),
			radial-gradient(ellipse 30% 6% at 45% 52%, rgb(255 244 225 / 0.4), transparent 70%);
		animation: drift 40s ease-in-out infinite alternate;
	}

	.grain {
		position: fixed;
		inset: 0;
		pointer-events: none;
		opacity: 0.18;
		mix-blend-mode: multiply;
		background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' width='160' height='160'><filter id='n'><feTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='2' stitchTiles='stitch'/><feColorMatrix values='0 0 0 0 0  0 0 0 0 0  0 0 0 0 0  0 0 0 0.7 0'/></filter><rect width='100%' height='100%' filter='url(%23n)'/></svg>");
	}

	.content {
		position: relative;
		max-width: 34rem;
		margin: 0 auto;
		padding: 2.5rem 1.25rem 3rem;
	}

	@keyframes drift {
		from {
			transform: translateX(-2%);
		}
		to {
			transform: translateX(2%);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.clouds {
			animation: none;
		}
	}
</style>
