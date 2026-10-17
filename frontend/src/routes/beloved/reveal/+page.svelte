<script lang="ts">
	import { onMount } from 'svelte';
	import { name, concert, note, password, matches, STORAGE_KEY } from '$lib/beloved/config';

	const unlockKey = STORAGE_KEY + ':unlocked';

	let unlocked = $state(false);
	let attempt = $state('');
	let error = $state('');
	let shake = $state(false);
	let ready = $state(false);

	const confetti = Array.from({ length: 36 }, (_, i) => ({
		left: (i * 97) % 100,
		delay: (i % 12) * 0.25,
		duration: 3.5 + (i % 5) * 0.6,
		color: ['#c4683a', '#e8923c', '#fbf1dd', '#6b3f5e'][i % 4]
	}));

	onMount(() => {
		unlocked = localStorage.getItem(unlockKey) === password;
		ready = true;
	});

	function submit(event: SubmitEvent) {
		event.preventDefault();
		if (matches(attempt, [password])) {
			localStorage.setItem(unlockKey, password);
			unlocked = true;
		} else {
			error = 'That is not it. The games will tell you.';
			shake = false;
			requestAnimationFrame(() => (shake = true));
		}
	}
</script>

<main class:ready>
	{#if unlocked}
		<div class="confetti" aria-hidden="true">
			{#each confetti as piece, i (i)}
				<span
					style="left:{piece.left}%;background:{piece.color};animation-delay:{piece.delay}s;animation-duration:{piece.duration}s"
				></span>
			{/each}
		</div>

		<p class="eyebrow">Happy birthday, {name}</p>
		<h1>We are going to see {concert.artist}</h1>

		<section class="ticket">
			<p class="tour">{concert.tour}</p>
			<p class="when">{concert.date}</p>
			<p class="where">{concert.venue}, {concert.city}</p>
			<p class="admit">Admit two</p>
		</section>

		<p class="note">{note}</p>
	{:else}
		<p class="eyebrow">One more door</p>
		<h1>{name}</h1>
		<form class="card" class:shake onsubmit={submit} onanimationend={() => (shake = false)}>
			<label for="pw">Enter the password you collected.</label>
			<input id="pw" bind:value={attempt} autocomplete="off" autocapitalize="off" />
			{#if error}<p class="error">{error}</p>{/if}
			<button class="btn" type="submit">Unlock</button>
			<a class="link" href="/beloved/">Back to the games</a>
		</form>
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
		font-size: clamp(2rem, 8vw, 3rem);
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
		text-align: center;
		color: var(--dusk);
	}

	.error {
		color: var(--clay);
		font-style: italic;
	}

	.shake {
		animation: shake 0.4s;
	}

	.ticket {
		position: relative;
		text-align: center;
		padding: 2rem 1.5rem;
		background: var(--cream);
		border-radius: 0.8rem;
		border: 2px dashed var(--clay);
		box-shadow: 0 12px 34px rgb(59 36 22 / 0.3);
		animation: rise 0.7s ease;
	}

	.tour {
		font-size: 1.5rem;
		font-style: italic;
		color: var(--clay);
	}

	.when {
		font-size: 1.9rem;
		font-weight: bold;
		margin: 0.4rem 0;
	}

	.where {
		font-size: 1.1rem;
	}

	.admit {
		margin-top: 1.2rem;
		padding-top: 0.8rem;
		border-top: 1px dashed var(--clay);
		letter-spacing: 0.25em;
		text-transform: uppercase;
		font-size: 0.8rem;
		color: var(--dusk);
	}

	.note {
		margin-top: 1.75rem;
		padding: 1.25rem;
		background: rgb(251 241 221 / 0.94);
		border-radius: 1rem;
		font-style: italic;
		white-space: pre-line;
	}

	.confetti {
		position: fixed;
		inset: 0;
		pointer-events: none;
		overflow: hidden;
	}

	.confetti span {
		position: absolute;
		top: -2rem;
		width: 0.6rem;
		height: 1rem;
		border-radius: 2px;
		animation: fall linear infinite;
	}

	@keyframes fall {
		to {
			transform: translateY(110vh) rotate(540deg);
		}
	}

	@keyframes rise {
		from {
			opacity: 0;
			transform: translateY(16px);
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
		.confetti {
			display: none;
		}

		.ticket,
		.shake {
			animation: none;
		}

		main {
			transition: none;
		}
	}
</style>
