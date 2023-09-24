export const ssr = false;
export const prerender = false;

import { getCollectionById, getItems } from "../../../../api";

export async function load({ fetch, params }) {
	let slug = params.slug;

	let res = await getCollectionById(fetch, slug);
	let resp = await res.json();
	let collection = resp.data;

	res = await getItems(fetch, slug);
	let items = await res.json();

	return {
		slug,
		collection,
		items
	};
}
