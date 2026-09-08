import { error } from '@sveltejs/kit';
import { cvEntries, getCvEntry } from '$lib/data/cv';
import type { EntryGenerator, PageLoad } from './$types';

// The timeline page only links work/education entries, so link-crawling alone
// would miss the pinned profile/skills notes. Enumerate them all explicitly.
export const entries: EntryGenerator = () => cvEntries.map((entry) => ({ slug: entry.slug }));

export const load: PageLoad = ({ params }) => {
	const entry = getCvEntry(params.slug);
	if (!entry) error(404, 'Not found');

	const { slug, ...metadata } = entry;
	return { slug, metadata };
};
