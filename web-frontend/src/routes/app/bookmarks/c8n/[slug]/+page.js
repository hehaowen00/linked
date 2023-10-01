export const ssr = false;
export const prerender = false;

import { redirect } from "@sveltejs/kit";
import { getCollectionById, getItems, loginUrl } from "../../../../../api";

export async function load({ fetch, url, params }) {
	try {
		let res = await validateUser(fetch, url.origin);
		if (!res.ok) {
			throw redirect(302, loginUrl(url.origin, url.href));
		}
	} catch (e) {}

	let slug = params.slug;

	let res = await getCollectionById(fetch, url.origin, slug);
	let resp = await res.json();
	let collection = resp.data;

	res = await getItems(fetch, url.origin, slug);
	let items = await res.json();

	return {
		slug,
		collection,
		items,
		url
	};
}
