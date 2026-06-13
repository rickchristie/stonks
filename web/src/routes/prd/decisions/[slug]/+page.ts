import { error } from '@sveltejs/kit';
import { getDecisionDoc } from '$lib/prd/data';

export function load({ params }) {
	const doc = getDecisionDoc(params.slug);
	if (!doc) throw error(404, 'Decision not found');
	return { doc };
}
