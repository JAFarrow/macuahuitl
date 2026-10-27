export const name = 'Micaela';

export interface Question {
	prompt: string;
	answers: string[];
	nudge: string;
}

export const questions: Question[] = [
	{
		prompt: 'PLACEHOLDER 1 — Where did we have our first date?',
		answers: ['placeholder one'],
		nudge: 'Close your eyes and think back to that evening…'
	},
	{
		prompt: 'PLACEHOLDER 2 — What is my favourite thing to order when we eat out?',
		answers: ['placeholder two'],
		nudge: 'You know this one. You have rolled your eyes at it often enough.'
	},
	{
		prompt: 'PLACEHOLDER 3 — Which song is unmistakably ours?',
		answers: ['placeholder three'],
		nudge: 'Hum it to yourself, the answer will follow.'
	}
];

export const password = 'loveyours';

export const shift = 4;

export const alphabet = 'abcdefghijklmnopqrstuvwxyz'.split('');

export function shifted(letter: string): string {
	return alphabet[(alphabet.indexOf(letter) + shift) % alphabet.length];
}

export const concert = {
	artist: 'J. Cole',
	tour: 'The Fall-Off Tour',
	date: '12 December 2026',
	venue: 'FNB Stadium',
	city: 'Johannesburg'
};

export function normalize(value: string): string {
	return value
		.normalize('NFD')
		.replace(/[\u0300-\u036f]/g, '')
		.toLowerCase()
		.replace(/[^a-z0-9]+/g, ' ')
		.trim();
}

export function matches(input: string, accepted: string[]): boolean {
	const value = normalize(input);
	return value !== '' && accepted.some((answer) => normalize(answer) === value);
}

export const STORAGE_KEY = 'beloved:v1';
