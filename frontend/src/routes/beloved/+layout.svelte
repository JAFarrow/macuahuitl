<script lang="ts">
	import type { Snippet } from 'svelte';
	import { name } from '#lib/beloved/config';

	let { children }: { children: Snippet } = $props();
</script>

<svelte:head>
	<title>For {name}</title>
	<meta name="robots" content="noindex, nofollow" />
	<meta name="theme-color" content="#0b0b0b" />
</svelte:head>

<div class="beloved">
	<div class="glow" aria-hidden="true"></div>
	<div class="grain" aria-hidden="true"></div>
	<div class="content">
		{@render children()}
	</div>
</div>

<style>
	:global(body) {
		margin: 0;
		background: #0b0b0b;
	}

	.beloved {
		--bg: #0b0b0b;
		--panel: #151412;
		--paper: #ece6d8;
		--ink: #0b0b0b;
		--dim: #a8a193;
		--accent: #ff5a1f;
		--display: Impact, Haettenschweiler, 'Arial Narrow Bold', 'Roboto Condensed', 'Arial Narrow',
			system-ui, sans-serif;
		--body: system-ui, -apple-system, 'Segoe UI', sans-serif;
		--mono: ui-monospace, 'SF Mono', 'Cascadia Mono', Menlo, Consolas, monospace;
		position: relative;
		box-sizing: border-box;
		min-height: 100vh;
		min-height: 100dvh;
		color: var(--paper);
		font-family: var(--body);
		line-height: 1.55;
		background: var(--bg);
		overflow-x: hidden;
	}

	.beloved :global(*),
	.beloved :global(*::before),
	.beloved :global(*::after) {
		box-sizing: border-box;
	}

	.beloved :global(:where(p, h1)) {
		margin: 0;
	}

	.glow {
		position: fixed;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		pointer-events: none;
		background:
			radial-gradient(ellipse 80% 45% at 50% 0%, rgb(255 90 31 / 0.22), transparent 70%),
			radial-gradient(ellipse 120% 60% at 50% 120%, rgb(0 0 0 / 0.9), transparent 70%);
		animation: breathe 8s ease-in-out infinite alternate;
	}

	.grain {
		position: fixed;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		pointer-events: none;
		opacity: 0.12;
		mix-blend-mode: screen;
		background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' width='160' height='160'><filter id='n'><feTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='2' stitchTiles='stitch'/><feColorMatrix values='1 0 0 0 0  1 0 0 0 0  1 0 0 0 0  0 0 0 0.6 0'/></filter><rect width='100%' height='100%' filter='url(%23n)'/></svg>");
	}

	.content {
		position: relative;
		max-width: 34rem;
		margin: 0 auto;
		padding: 2.5rem 1.25rem 3rem;
	}

	.beloved :global(.eyebrow) {
		text-align: center;
		margin-top: 1.5rem;
		font-size: 1.3rem;
		font-style: italic;
		letter-spacing: 0.02em;
		color: var(--accent);
	}

	.beloved :global(h1) {
		text-align: center;
		margin-bottom: 1.75rem;
		font-family: var(--display);
		font-weight: 400;
		font-size: clamp(2.6rem, 13vw, 4.2rem);
		line-height: 0.95;
		letter-spacing: 0.01em;
		text-transform: uppercase;
		color: var(--paper);
	}

	.beloved :global(.card) {
		display: flex;
		flex-direction: column;
		padding: 1.5rem;
		background: var(--panel);
		border: 2px solid var(--paper);
		box-shadow: 8px 8px 0 var(--accent);
	}

	.beloved :global(.card > * + *) {
		margin-top: 1rem;
	}

	.beloved :global(.card p) {
		color: var(--paper);
	}

	.beloved :global(.step) {
		font-family: var(--mono);
		font-size: 0.75rem;
		letter-spacing: 0.25em;
		text-transform: uppercase;
		color: var(--accent);
	}

	.beloved :global(label) {
		font-size: 1.1rem;
		font-weight: 600;
	}

	.beloved :global(input) {
		font: inherit;
		font-family: var(--mono);
		padding: 0.75rem 0.9rem;
		border: 2px solid var(--paper);
		border-radius: 0;
		background: var(--bg);
		color: var(--paper);
	}

	.beloved :global(input:focus-visible),
	.beloved :global(.btn:focus-visible),
	.beloved :global(.link:focus-visible) {
		outline: 3px solid var(--accent);
		outline-offset: 3px;
	}

	.beloved :global(.btn) {
		display: inline-block;
		text-align: center;
		text-decoration: none;
		font-family: var(--display);
		font-size: 1.25rem;
		letter-spacing: 0.12em;
		text-transform: uppercase;
		padding: 0.7rem 1.4rem;
		border: 0;
		border-radius: 0;
		background: var(--accent);
		color: var(--ink);
		box-shadow: 4px 4px 0 var(--paper);
		cursor: pointer;
		transition:
			transform 0.1s ease,
			box-shadow 0.1s ease;
	}

	.beloved :global(.btn:hover) {
		transform: translate(2px, 2px);
		box-shadow: 2px 2px 0 var(--paper);
	}

	.beloved :global(.btn:active) {
		transform: translate(4px, 4px);
		box-shadow: 0 0 0 var(--paper);
	}

	.beloved :global(.link) {
		font: inherit;
		font-family: var(--mono);
		font-size: 0.85rem;
		letter-spacing: 0.08em;
		text-align: center;
		min-height: 2.75rem;
		padding: 0.6rem 0.75rem;
		background: none;
		border: 0;
		color: var(--dim);
		text-decoration: underline;
		text-underline-offset: 0.25em;
		cursor: pointer;
	}

	.beloved :global(.link:hover) {
		color: var(--paper);
	}

	.beloved :global(.nudge),
	.beloved :global(.error) {
		font-family: var(--mono);
		font-size: 0.9rem;
		color: var(--accent);
	}

	.beloved :global(.shake) {
		animation: shake 0.4s;
	}

	@keyframes shake {
		25% {
			transform: translateX(-6px);
		}
		75% {
			transform: translateX(6px);
		}
	}

	@keyframes breathe {
		from {
			opacity: 0.7;
		}
		to {
			opacity: 1;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.glow,
		.beloved :global(.shake),
		.beloved :global(.btn) {
			animation: none;
			transition: none;
		}
	}
</style>
