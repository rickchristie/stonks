import { error } from '@sveltejs/kit';
import { getFeatureDoc } from '$lib/prd/data';

export function load({ params }) {
	const doc = getFeatureDoc(params.slug);
	if (!doc) throw error(404, 'Feature not found');
	return { doc };
}
