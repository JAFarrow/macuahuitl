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
				Happy early birthday. There is a surprise waiting for you, but it is locked behind a
				password, and the password is written in a letter that only a secret key can open.
			</p>
			<p>Answer a few questions to earn the key.</p>
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
					<span>{praise}</span>
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

	.stamp {
		position: absolute;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		display: grid;
		place-items: center;
		background: rgb(21 20 18 / 0.9);
		animation:
			fade 0.6s ease both,
			leave 0.6s ease 2.4s forwards;
	}

	.stamp span {
		font-family: var(--display);
		font-size: clamp(2.2rem, 9vw, 3.6rem);
		letter-spacing: 0.04em;
		white-space: nowrap;
		text-transform: uppercase;
		color: var(--ink);
		background: var(--accent);
		padding: 0.1em 0.5em;
		box-shadow: 6px 6px 0 var(--paper);
		transform: rotate(-5deg);
		animation: stamp 1.1s cubic-bezier(0.3, 0.9, 0.3, 1) 0.2s both;
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
		border: 2px solid var(--paper);
		box-shadow: 5px 5px 0 var(--accent);
		-webkit-backface-visibility: hidden;
		backface-visibility: hidden;
		font-family: var(--display);
		font-size: clamp(2.2rem, 11vw, 3.8rem);
		text-transform: uppercase;
	}

	.plain {
		color: var(--paper);
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
		0% {
			box-shadow: 8px 8px 0 var(--accent);
		}
		30% {
			transform: scale(1.03);
			box-shadow: 0 0 0 6px var(--accent);
		}
		100% {
			box-shadow: 8px 8px 0 var(--accent);
		}
	}

	@keyframes stamp {
		0% {
			opacity: 0;
			transform: rotate(-18deg) scale(2.6);
		}
		45% {
			opacity: 1;
			transform: rotate(-4deg) scale(0.94);
		}
		70% {
			transform: rotate(-6deg) scale(1.03);
		}
		100% {
			transform: rotate(-5deg) scale(1);
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
