<script lang="ts">
	import { onMount } from 'svelte';
	import { name, questions, alphabet, shifted, matches, STORAGE_KEY } from '$lib/beloved/config';

	const doneStep = questions.length + 1;
	const sample = alphabet.slice(0, 5);

	let step = $state(0);
	let answer = $state('');
	let misses = $state(0);
	let nudge = $state('');
	let shake = $state(false);
	let ready = $state(false);
	let run = $state(0);

	const question = $derived(step >= 1 && step <= questions.length ? questions[step - 1] : null);

	onMount(() => {
		const saved = Number(localStorage.getItem(STORAGE_KEY + ':step'));
		if (Number.isInteger(saved) && saved >= 0 && saved <= doneStep) step = saved;
		ready = true;
	});

	function go(next: number) {
		step = next;
		answer = '';
		misses = 0;
		nudge = '';
		localStorage.setItem(STORAGE_KEY + ':step', String(next));
	}

	function submit(event: SubmitEvent) {
		event.preventDefault();
		if (!question) return;
		if (matches(answer, question.answers)) {
			go(step + 1);
		} else {
			misses += 1;
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
		<form class="card" class:shake onsubmit={submit} onanimationend={() => (shake = false)}>
			<p class="step">Question {step} of {questions.length}</p>
			<label for="answer">{question.prompt}</label>
			<input id="answer" bind:value={answer} autocomplete="off" autocapitalize="off" />
			{#if nudge}<p class="nudge">{nudge}</p>{/if}
			<button class="btn" type="submit">Check</button>
			{#if misses >= 3}
				<button class="link" type="button" onclick={() => go(step + 1)}>Stuck? Skip this one</button>
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
		animation: rise 0.45s ease;
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

	@media (prefers-reduced-motion: reduce) {
		.card,
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
