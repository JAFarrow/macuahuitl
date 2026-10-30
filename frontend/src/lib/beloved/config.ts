export const name = 'Micaela';

export interface Question {
	prompt: string;
	answers: string[];
	nudge: string;
}

export const questions: Question[] = [
	{
		prompt: 'With a dream in mind, the sea in front of you, turn right, across the garden, you find...',
		answers: ['a studio', 'studio', 'art studio', 'my studio'],
		nudge: 'Where childrens\' books come to life.'
	},
	{
		prompt: 'For most, an expression. For us, a pet name.',
		answers: ['scrunch'],
		nudge: 'When you\'re feeling particularly shy.'
	},
	{
		prompt: 'Fluttering, and kept. Full, and quiet.',
		answers: [
			'a cageful of butterflies',
			'cageful of butterflies',
			'cage full of butterflies',
			'a cage full of butterflies',
			'a cage of butterflies',
			'cage of butterflies'
		],
		nudge: 'It\'s the one I read to you when the day has been too much.'
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
