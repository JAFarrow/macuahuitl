<script lang="ts">
	import type { Snippet } from 'svelte';
	import { name } from '#lib/beloved/config';

	let { children }: { children: Snippet } = $props();

	const butterflies = [
		{ left: 6, top: 12, size: 56, delay: 0, duration: 19, color: '#6f94c0' },
		{ left: 80, top: 20, size: 40, delay: -6, duration: 23, color: '#e08f8a' },
		{ left: 12, top: 66, size: 36, delay: -11, duration: 21, color: '#d9a441' },
		{ left: 74, top: 72, size: 64, delay: -3, duration: 26, color: '#8fae8b' },
		{ left: 48, top: 4, size: 30, delay: -15, duration: 17, color: '#e08f8a' }
	];
</script>

<svelte:head>
	<title>For {name}</title>
	<meta name="robots" content="noindex, nofollow" />
	<meta name="theme-color" content="#f6efe2" />
</svelte:head>

<div class="beloved">
	<div class="glow" aria-hidden="true"></div>
	<div class="grain" aria-hidden="true"></div>
	<div class="flutter" aria-hidden="true">
		{#each butterflies as b, i (i)}
			<span
				style="left:{b.left}%;top:{b.top}%;width:{b.size}px;color:{b.color};animation-delay:{b.delay}s;animation-duration:{b.duration}s"
			>
				<svg viewBox="0 0 64 48">
					<g fill="currentColor">
						<path d="M32 24C22 4 4 2 3 14C2 24 14 30 32 24Z" />
						<path d="M32 24C42 4 60 2 61 14C62 24 50 30 32 24Z" />
						<path d="M32 24C20 28 10 44 20 45C28 46 32 34 32 24Z" opacity="0.8" />
						<path d="M32 24C44 28 54 44 44 45C36 46 32 34 32 24Z" opacity="0.8" />
					</g>
					<rect x="31" y="14" width="2" height="22" rx="1" fill="#3b2f2a" opacity="0.6" />
				</svg>
			</span>
		{/each}
	</div>
	<div class="content">
		{@render children()}
	</div>
</div>

<style>
	:global(body) {
		margin: 0;
		background: #f6efe2;
	}

	.beloved {
		--bg: #f6efe2;
		--panel: #fffaf0;
		--ticket: #fffaf0;
		--text: #3b2f2a;
		--on-accent: #fffaf0;
		--dim: #6f6156;
		--line: #cdbb9c;
		--accent: #3f6590;
		--shadow: 0 14px 36px rgb(59 47 42 / 0.22);
		--display: 'Iowan Old Style', 'Palatino Linotype', Palatino, Georgia, 'Times New Roman', serif;
		--body: system-ui, -apple-system, 'Segoe UI', sans-serif;
		--mono: ui-monospace, 'SF Mono', 'Cascadia Mono', Menlo, Consolas, monospace;
		position: relative;
		box-sizing: border-box;
		min-height: 100vh;
		min-height: 100dvh;
		color: var(--text);
		font-family: var(--body);
		line-height: 1.6;
		background: linear-gradient(180deg, #fbf5e9 0%, var(--bg) 35%, #e6d2ab 100%);
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
			radial-gradient(ellipse 80% 45% at 50% 0%, rgb(232 180 160 / 0.35), transparent 70%),
			radial-gradient(ellipse 120% 60% at 50% 120%, rgb(217 164 65 / 0.18), transparent 70%);
		animation: breathe 8s ease-in-out infinite alternate;
	}

	.grain {
		position: fixed;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		pointer-events: none;
		opacity: 0.05;
		mix-blend-mode: multiply;
		background-image: url("data:image/svg+xml;utf8,<svg xmlns='http://www.w3.org/2000/svg' width='160' height='160'><filter id='n'><feTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='2' stitchTiles='stitch'/><feColorMatrix values='1 0 0 0 0  1 0 0 0 0  1 0 0 0 0  0 0 0 0.6 0'/></filter><rect width='100%' height='100%' filter='url(%23n)'/></svg>");
	}

	.flutter {
		position: fixed;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		pointer-events: none;
		overflow: hidden;
	}

	.flutter span {
		position: absolute;
		opacity: 0.75;
		animation: drift ease-in-out infinite alternate;
	}

	.flutter svg {
		display: block;
		width: 100%;
		height: auto;
		animation: flap 1.4s ease-in-out infinite alternate;
	}

	.content {
		position: relative;
		max-width: 34rem;
		margin: 0 auto;
		padding: 1.25rem 1.25rem 1.5rem;
	}

	.beloved :global(.eyebrow) {
		text-align: center;
		margin-top: 0.75rem;
		font-family: var(--display);
		font-size: 1.3rem;
		font-style: italic;
		letter-spacing: 0.02em;
		color: var(--accent);
	}

	.beloved :global(h1) {
		text-align: center;
		margin-bottom: 1.1rem;
		font-family: var(--display);
		font-weight: 400;
		font-size: clamp(2.4rem, 12vw, 4.2rem);
		line-height: 1.05;
		letter-spacing: 0.005em;
		text-wrap: balance;
		color: var(--text);
	}

	.beloved :global(h1)::after {
		content: '';
		display: block;
		width: 9rem;
		height: 9px;
		margin: 0.6rem auto 0;
		background:
			radial-gradient(circle at center, var(--accent) 3px, transparent 4px),
			linear-gradient(90deg, transparent, var(--accent) 35%, var(--accent) 65%, transparent) center /
				100% 1px no-repeat;
		opacity: 0.7;
	}

	.beloved :global(.card) {
		display: flex;
		flex-direction: column;
		padding: 1.25rem 1.25rem;
		background: linear-gradient(180deg, #fffdf6, var(--panel) 40%, #fbf3e2);
		border: 1px solid var(--line);
		border-radius: 14px;
		box-shadow:
			var(--shadow),
			0 2px 0 rgb(255 255 255 / 0.8) inset,
			0 0 0 6px rgb(255 250 240 / 0.55);
	}

	.beloved :global(.card > * + *) {
		margin-top: 0.8rem;
	}

	.beloved :global(.card p) {
		color: var(--text);
	}

	.beloved :global(.card p:not(.step, .nudge, .error)) {
		font-family: var(--display);
		font-size: 1.15rem;
		line-height: 1.7;
	}

	.beloved :global(.step) {
		font-family: var(--display);
		font-size: 0.95rem;
		font-style: italic;
		letter-spacing: 0.04em;
		color: var(--accent);
	}

	.beloved :global(label) {
		font-family: var(--display);
		font-size: 1.2rem;
		line-height: 1.45;
	}

	.beloved :global(input) {
		font: inherit;
		font-family: var(--display);
		padding: 0.75rem 0.9rem;
		border: 1px solid var(--line);
		border-radius: 10px;
		background: var(--bg);
		color: var(--text);
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
		font-size: 1.1rem;
		letter-spacing: 0.04em;
		padding: 0.7rem 1.6rem;
		border: 0;
		border-radius: 999px;
		background: var(--accent);
		color: var(--on-accent);
		box-shadow: 0 4px 12px rgb(63 101 144 / 0.25);
		cursor: pointer;
		transition:
			transform 0.15s ease,
			box-shadow 0.15s ease;
	}

	.beloved :global(.btn:hover) {
		transform: translateY(-1px);
		box-shadow: 0 6px 16px rgb(63 101 144 / 0.3);
	}

	.beloved :global(.btn:active) {
		transform: translateY(1px);
		box-shadow: 0 2px 6px rgb(63 101 144 / 0.25);
	}

	.beloved :global(.link) {
		font: inherit;
		font-family: var(--display);
		font-size: 0.95rem;
		font-style: italic;
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
		color: var(--text);
	}

	.beloved :global(.nudge),
	.beloved :global(.error) {
		font-family: var(--display);
		font-size: 1rem;
		font-style: italic;
		color: var(--accent);
	}

	.beloved :global(.shake) {
		animation: shake 0.4s;
	}

	@keyframes shake {
		25% {
			transform: translateX(-4px);
		}
		75% {
			transform: translateX(4px);
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

	@keyframes drift {
		from {
			transform: translate(0, 0) rotate(-8deg);
		}
		to {
			transform: translate(28px, -46px) rotate(10deg);
		}
	}

	@keyframes flap {
		from {
			transform: scaleX(1);
		}
		to {
			transform: scaleX(0.78);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.flutter span,
		.flutter svg,
		.glow,
		.beloved :global(.shake),
		.beloved :global(.btn) {
			animation: none;
			transition: none;
		}
	}
</style>
