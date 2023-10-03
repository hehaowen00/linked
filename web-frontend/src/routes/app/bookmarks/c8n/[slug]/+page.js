export const ssr = false;
export const prerender = false;

import { getCollectionById, getItems } from "$lib/api";

export async function load({ fetch, url, params }) {
	let slug = params.slug;

	let res = await getCollectionById(fetch, url.origin, slug);
	let resp = await res.json();
	let collection = resp.data;

	res = await getItems(fetch, url.origin, slug);
	let items = await res.json();

	return {
		slug,
		collection,
		items
	};
}
