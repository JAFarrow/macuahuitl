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
	<p class="eyebrow">A little something for</p>
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
							<span class="glyph plain">{letter}</span>
							<span class="glyph coded">{shifted(letter)}</span>
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

	.eyebrow {
		text-align: center;
		letter-spacing: 0.2em;
		text-transform: uppercase;
		font-size: 0.8rem;
		color: var(--cream);
		margin-top: 1.5rem;
	}

	h1 {
		text-align: center;
		font-size: clamp(2.2rem, 9vw, 3.2rem);
		font-weight: 400;
		font-style: italic;
		color: var(--cream);
		margin-bottom: 1.5rem;
		text-shadow: 0 2px 12px rgb(59 36 22 / 0.35);
	}

	.card {
		display: flex;
		flex-direction: column;
		gap: 0.9rem;
		padding: 1.5rem;
		background: rgb(251 241 221 / 0.94);
		border-radius: 1rem;
		box-shadow: 0 10px 30px rgb(59 36 22 / 0.25);
		animation: rise 0.45s ease;
	}

	.step {
		font-size: 0.8rem;
		letter-spacing: 0.15em;
		text-transform: uppercase;
		color: var(--clay);
	}

	label {
		font-size: 1.15rem;
	}

	input {
		font: inherit;
		padding: 0.7rem 0.9rem;
		border: 2px solid var(--ink);
		border-radius: 0.6rem;
		background: #fff;
		color: var(--ink);
	}

	input:focus-visible,
	.btn:focus-visible,
	.link:focus-visible {
		outline: 3px solid var(--sun);
		outline-offset: 2px;
	}

	.btn {
		display: inline-block;
		text-align: center;
		text-decoration: none;
		font: inherit;
		padding: 0.75rem 1.2rem;
		border: 0;
		border-radius: 999px;
		background: var(--ink);
		color: var(--cream);
		cursor: pointer;
	}

	.btn:hover {
		background: var(--clay);
	}

	.link {
		font: inherit;
		background: none;
		border: 0;
		color: var(--dusk);
		text-decoration: underline;
		cursor: pointer;
	}

	.nudge {
		color: var(--clay);
		font-style: italic;
	}

	.reveal {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 2.5rem;
		padding-top: 1.5rem;
	}

	.cipher {
		display: flex;
		justify-content: center;
		gap: clamp(0.5rem, 3vw, 1rem);
		list-style: none;
		padding: 0;
	}

	.cipher li {
		position: relative;
		width: clamp(3.6rem, 15vw, 4.6rem);
		height: clamp(5rem, 21vw, 6.4rem);
		border-radius: 0.9rem;
		background: var(--cream);
		box-shadow: 0 10px 26px rgb(59 36 22 / 0.3);
		overflow: hidden;
		animation: rise 0.6s ease both;
		animation-delay: calc(var(--i) * 0.18s);
	}

	.glyph {
		position: absolute;
		inset: 0;
		display: grid;
		place-items: center;
		font-size: clamp(2.4rem, 10vw, 3.2rem);
		font-style: italic;
	}

	.plain {
		color: var(--ink);
		animation: out 0.55s ease both;
		animation-delay: calc(1.3s + var(--i) * 0.45s);
	}

	.coded {
		color: var(--clay);
		font-weight: bold;
		animation: in 0.55s ease both;
		animation-delay: calc(1.3s + var(--i) * 0.45s);
	}

	.actions {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 0.9rem;
		animation: rise 0.6s ease both;
		animation-delay: 3.8s;
	}

	.actions .link {
		color: var(--cream);
	}

	.shake {
		animation: shake 0.4s;
	}

	@keyframes rise {
		from {
			opacity: 0;
			transform: translateY(12px);
		}
	}

	@keyframes out {
		to {
			transform: translateY(-110%);
			opacity: 0;
		}
	}

	@keyframes in {
		from {
			transform: translateY(110%);
			opacity: 0;
		}
	}

	@keyframes shake {
		25% {
			transform: translateX(-6px);
		}
		75% {
			transform: translateX(6px);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.card,
		.shake,
		.cipher li,
		.actions,
		.plain,
		.coded {
			animation: none;
		}

		.plain {
			display: none;
		}

		main {
			transition: none;
		}
	}
</style>
