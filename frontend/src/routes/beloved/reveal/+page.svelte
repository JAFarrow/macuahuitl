<script lang="ts">
	import { onMount } from 'svelte';
	import { name, concert, password, matches, STORAGE_KEY } from '#lib/beloved/config';

	const unlockKey = STORAGE_KEY + ':unlocked';

	let unlocked = $state(false);
	let attempt = $state('');
	let error = $state('');
	let shake = $state(false);
	let ready = $state(false);

	const confetti = Array.from({ length: 36 }, (_, i) => ({
		left: (i * 97) % 100,
		delay: (i % 12) * 0.25,
		duration: 4.5 + (i % 5) * 0.7,
		color: ['#7f9fc4', '#e3a6a0', '#d9a441', '#9db39a'][i % 4]
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
			error = 'Not quite. Hint: it is one of your favourite songs.';
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
		<h1>We are going to see {concert.artist.replace(/ /g, ' ')}!</h1>

		<section class="ticket">
			<div class="top">
				<p class="tour">{concert.tour}</p>
				<p class="when">{concert.date}</p>
				<p class="where">{concert.venue}, {concert.city}</p>
			</div>
			<div class="admit">
				<p>Admit two</p>
				<div class="barcode" aria-hidden="true"></div>
			</div>
		</section>
	{:else}
		<p class="eyebrow">One more door</p>
		<h1>{name}</h1>
		<form class="card" class:shake onsubmit={submit} onanimationend={() => (shake = false)}>
			<label for="pw">Enter the password you collected.</label>
			<input id="pw" bind:value={attempt} autocomplete="off" autocapitalize="off" />
			{#if error}<p class="error">{error}</p>{/if}
			<button class="btn" type="submit">Unlock</button>
			<a class="link" href="/beloved/">Back to the key</a>
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

	main :global(.eyebrow),
	main :global(h1),
	.ticket {
		position: relative;
		z-index: 1;
	}

	:global(.beloved .content) main h1 {
		margin-top: 0.4rem;
		margin-bottom: 1.25rem;
		font-size: clamp(2rem, 9.5vw, 3.6rem);
		line-height: 1.1;
	}

	.ticket {
		text-align: center;
		color: var(--text);
		filter: drop-shadow(0 14px 22px rgb(59 47 42 / 0.3)) drop-shadow(0 0 1px rgb(59 47 42 / 0.5));
		animation: rise 0.7s ease;
	}

	.top,
	.admit {
		background: var(--ticket);
	}

	.top {
		padding: 1.5rem 1.25rem 1.25rem;
		-webkit-mask:
			radial-gradient(circle 8px at 0 100%, transparent 96%, #000 100%),
			radial-gradient(circle 8px at 100% 100%, transparent 96%, #000 100%);
		-webkit-mask-composite: source-in;
		mask:
			radial-gradient(circle 8px at 0 100%, transparent 96%, #000 100%),
			radial-gradient(circle 8px at 100% 100%, transparent 96%, #000 100%);
		mask-composite: intersect;
	}

	.ticket p {
		color: var(--text);
	}

	.tour {
		font-family: var(--display);
		font-size: 1rem;
		font-style: italic;
		letter-spacing: 0.06em;
	}

	.when {
		margin: 0.6rem 0 0.3rem;
		font-family: var(--display);
		font-size: clamp(2rem, 9vw, 3.2rem);
		line-height: 1.05;
		color: var(--accent);
	}

	.ticket .when {
		color: var(--accent);
	}

	.where {
		font-family: var(--display);
		font-size: 1.15rem;
		letter-spacing: 0.03em;
	}

	.admit {
		padding: 0.9rem 1.25rem 1.1rem;
		border-top: 2px dashed #a8946f;
		font-family: var(--display);
		font-size: 1rem;
		font-style: italic;
		letter-spacing: 0.08em;
		-webkit-mask:
			radial-gradient(circle 8px at 0 0, transparent 96%, #000 100%),
			radial-gradient(circle 8px at 100% 0, transparent 96%, #000 100%);
		-webkit-mask-composite: source-in;
		mask:
			radial-gradient(circle 8px at 0 0, transparent 96%, #000 100%),
			radial-gradient(circle 8px at 100% 0, transparent 96%, #000 100%);
		mask-composite: intersect;
	}

	.barcode {
		height: 2.25rem;
		margin-top: 0.7rem;
		opacity: 0.45;
		background: repeating-linear-gradient(
			90deg,
			var(--text) 0 2px,
			transparent 2px 5px,
			var(--text) 5px 6px,
			transparent 6px 10px,
			var(--text) 10px 13px,
			transparent 13px 15px
		);
	}

	.confetti {
		position: fixed;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		pointer-events: none;
		overflow: hidden;
	}

	.confetti span {
		position: absolute;
		top: -2rem;
		width: 0.7rem;
		height: 0.7rem;
		border-radius: 0 100% 0 100%;
		opacity: 0.7;
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

	@media (prefers-reduced-motion: reduce) {
		.confetti {
			display: none;
		}

		.ticket {
			animation: none;
		}

		main {
			transition: none;
		}
	}
</style>
