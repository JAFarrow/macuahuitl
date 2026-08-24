import { error } from '@sveltejs/kit';
import { getPost } from '$lib/data/posts';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ params }) => {
	const post = getPost(params.slug);
	if (!post) error(404, 'Not found');

	const { slug, ...metadata } = post;
	return { slug, metadata };
};
