<script lang="ts">
	import { onMount } from 'svelte';
	import { name, questions, alphabet, shifted, matches, STORAGE_KEY } from '#lib/beloved/config';

	const doneStep = questions.length + 1;
	const sample = alphabet.slice(0, 5);
	const praises = ['Correct', "That's it", 'Right'];
	let step = $state(0);
	let answer = $state('');
	let nudge = $state('');
	let shake = $state(false);
	let ready = $state(false);
	let run = $state(0);
	let cheer = $state(false);
	let praise = $state('');
	let timer: ReturnType<typeof setTimeout> | undefined;

	const question = $derived(step >= 1 && step <= questions.length ? questions[step - 1] : null);

	onMount(() => {
		const saved = Number(localStorage.getItem(STORAGE_KEY + ':step'));
		if (Number.isInteger(saved) && saved >= 0 && saved <= doneStep) step = saved;
		ready = true;
		return () => clearTimeout(timer);
	});

	function celebrate() {
		const last = step === questions.length;
		const calm = matchMedia('(prefers-reduced-motion: reduce)').matches;
		praise = last ? 'You\'re in!' : praises[(step - 1) % praises.length];
		cheer = true;
		navigator.vibrate?.([40, 40, 90]);
		const next = step + 1;
		localStorage.setItem(STORAGE_KEY + ':step', String(next));
		timer = setTimeout(
			() => {
				cheer = false;
				go(next);
			},
			calm ? 800 : 3000
		);
	}

	function go(next: number) {
		step = next;
		answer = '';
		nudge = '';
		localStorage.setItem(STORAGE_KEY + ':step', String(next));
	}

	function submit(event: SubmitEvent) {
		event.preventDefault();
		if (!question || cheer) return;
		if (matches(answer, question.answers)) {
			celebrate();
		} else {
			nudge = question.nudge;
			shake = false;
			requestAnimationFrame(() => (shake = true));
		}
	}
</script>

<main class:ready>
	<p class="eyebrow">Just for you,</p>
	<h1>{name}</h1>

	{#if step === 0}
		<section class="card">
			<p>
				Happy birthday. I've kept something for you, and it only opens for you. Answer me three
				things and it's yours.
			</p>
			<button class="btn" onclick={() => go(1)}>Begin</button>
		</section>
	{:else if question}
		<form
			class="card"
			class:shake
			class:cheer
			onsubmit={submit}
			onanimationend={(e) => {
				if (e.animationName === 'shake') shake = false;
			}}
		>
			<p class="step">Question {step} of {questions.length}</p>
			<label for="answer">{question.prompt}</label>
			<input
				id="answer"
				bind:value={answer}
				disabled={cheer}
				autocomplete="off"
				autocapitalize="off"
			/>
			{#if nudge}<p class="nudge">{nudge}</p>{/if}
			<button class="btn" type="submit" disabled={cheer}>Check</button>
			{#if cheer}
				<div class="stamp" role="status">
					<svg class="moth" viewBox="0 0 64 48" aria-hidden="true">
						<g class="wings" fill="currentColor">
							<path d="M32 24C22 4 4 2 3 14C2 24 14 30 32 24Z" />
							<path d="M32 24C42 4 60 2 61 14C62 24 50 30 32 24Z" />
							<path d="M32 24C20 28 10 44 20 45C28 46 32 34 32 24Z" opacity="0.8" />
							<path d="M32 24C44 28 54 44 44 45C36 46 32 34 32 24Z" opacity="0.8" />
						</g>
						<rect x="31" y="14" width="2" height="22" rx="1" fill="#3b2f2a" opacity="0.6" />
					</svg>
					<span>{praise}</span>
					<b class="rule" aria-hidden="true"></b>
				</div>
			{/if}
		</form>
	{:else}
		{#key run}
			<section class="reveal" aria-label="Cipher key">
				<ul class="cipher">
					{#each sample as letter, i (letter)}
						<li style="--i:{i}">
							<span class="flip">
								<span class="glyph plain">{letter}</span>
								<span class="glyph coded">{shifted(letter)}</span>
							</span>
						</li>
					{/each}
				</ul>
				<div class="actions">
					<a class="btn" href="/beloved/reveal/">Open the surprise</a>
					<button class="link" type="button" onclick={() => run++}>Replay</button>
				</div>
			</section>
		{/key}
	{/if}
</main>

<style>
	main {
		opacity: 0;
		transition: opacity 0.4s ease;
	}

	main.ready {
		opacity: 1;
	}

	.card {
		position: relative;
		animation: rise 0.45s ease;
	}

	.card.cheer {
		animation: pulse 0.9s ease;
	}

	.card > .stamp {
		margin: 0;
		position: absolute;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		overflow: hidden;
		border-radius: 14px;
		background: rgb(246 239 226 / 0.94);
		animation:
			fade 0.6s ease both,
			leave 0.6s ease 2.4s forwards;
	}

	.stamp::before {
		content: '';
		position: absolute;
		top: 50%;
		left: 50%;
		width: 120%;
		aspect-ratio: 1;
		margin: -60% 0 0 -60%;
		background: radial-gradient(circle, rgb(217 164 65 / 0.35), rgb(63 101 144 / 0.1) 40%, transparent 65%);
		animation: glow 1.6s ease-out 0.2s both;
	}

	.moth {
		position: relative;
		width: 3.6rem;
		color: var(--accent);
		animation: moth 1.7s cubic-bezier(0.3, 0.5, 0.3, 1) 0.1s both;
	}

	.wings {
		transform-box: fill-box;
		transform-origin: center;
		animation: flutter 0.3s ease-in-out 0.1s 6 alternate;
	}

	.stamp span {
		position: relative;
		margin-top: 0.4rem;
		font-family: var(--display);
		font-size: clamp(2.2rem, 9vw, 3.4rem);
		font-style: italic;
		letter-spacing: 0.01em;
		white-space: nowrap;
		color: var(--accent);
		animation: praise 0.8s cubic-bezier(0.2, 0.9, 0.3, 1.2) 1.1s both;
	}

	.rule {
		position: relative;
		width: 7rem;
		height: 2px;
		margin-top: 0.5rem;
		background: linear-gradient(90deg, transparent, #d9a441, transparent);
		transform-origin: center;
		animation: draw 0.6s ease 1.6s both;
	}

	.reveal {
		display: flex;
		flex-direction: column;
		align-items: center;
		padding-top: 1.5rem;
	}

	.reveal > :global(* + *) {
		margin-top: 2.5rem;
	}

	.cipher {
		display: flex;
		justify-content: center;
		list-style: none;
		margin: 0;
		padding: 0;
	}

	.cipher li + li {
		margin-left: clamp(0.5rem, 3vw, 0.8rem);
	}

	.cipher li {
		position: relative;
		width: clamp(2.8rem, 14vw, 4.6rem);
		height: clamp(4.4rem, 19vw, 6.4rem);
		perspective: 600px;
		animation: rise 0.6s ease both;
		animation-delay: calc(var(--i) * 0.18s);
	}

	.flip {
		position: absolute;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		transform-style: preserve-3d;
		animation: flip 0.8s cubic-bezier(0.45, 0, 0.2, 1) both;
		animation-delay: calc(1.3s + var(--i) * 0.45s);
	}

	.glyph {
		position: absolute;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		display: grid;
		place-items: center;
		background: var(--panel);
		border: 1px solid var(--line);
		border-radius: 12px;
		box-shadow: var(--shadow);
		-webkit-backface-visibility: hidden;
		backface-visibility: hidden;
		font-family: var(--display);
		font-size: clamp(2.2rem, 11vw, 3.8rem);
		font-style: italic;
	}

	.plain {
		color: var(--text);
	}

	.coded {
		color: var(--accent);
		transform: rotateX(180deg);
	}

	.actions {
		display: flex;
		flex-direction: column;
		align-items: center;
		animation: rise 0.6s ease both;
		animation-delay: 3.8s;
	}

	.actions > :global(* + *) {
		margin-top: 1rem;
	}

	@keyframes rise {
		from {
			opacity: 0;
			transform: translateY(12px);
		}
	}

	@keyframes flip {
		to {
			transform: rotateX(180deg);
		}
	}

	@keyframes pulse {
		0%,
		100% {
			box-shadow: var(--shadow);
		}
		40% {
			box-shadow:
				var(--shadow),
				0 0 0 8px rgb(63 101 144 / 0.15);
		}
	}

	@keyframes praise {
		from {
			opacity: 0;
			transform: translateY(10px) scale(0.85);
		}
	}

	@keyframes moth {
		from {
			opacity: 0;
			transform: translate(-90px, 70px) rotate(-35deg) scale(0.3);
		}
		60% {
			opacity: 1;
			transform: translate(10px, -8px) rotate(8deg) scale(1.1);
		}
	}

	@keyframes flutter {
		from {
			transform: scaleX(1);
		}
		to {
			transform: scaleX(0.4);
		}
	}

	@keyframes glow {
		from {
			opacity: 0;
			transform: scale(0.3);
		}
		40% {
			opacity: 1;
		}
		to {
			opacity: 0.55;
			transform: scale(1);
		}
	}

	@keyframes draw {
		from {
			opacity: 0;
			transform: scaleX(0);
		}
	}

	@keyframes leave {
		to {
			opacity: 0;
		}
	}

	@keyframes fade {
		from {
			opacity: 0;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.card,
		.card.cheer,
		.stamp,
		.stamp::before,
		.moth,
		.wings,
		.rule,
		.stamp span,
		.cipher li,
		.flip,
		.actions {
			animation: none;
		}

		.flip {
			transform: rotateX(180deg);
		}

		main {
			transition: none;
		}
	}
</style>
