import { error } from '@sveltejs/kit';
import { getProject } from '$lib/data/projects';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ params }) => {
	const project = getProject(params.slug);
	if (!project) error(404, 'Not found');

	const { slug, ...metadata } = project;
	return { slug, metadata };
};
