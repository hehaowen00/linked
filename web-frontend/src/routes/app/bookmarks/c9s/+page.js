import { getCollections } from "$lib/api";

export const ssr = false;

export async function load({ fetch, url }) {
	let res = await getCollections(fetch, url.origin);
	if (!res.ok) {
		return;
	}
	let json = await res.json();
	return {
		collections: json.data ?? []
	};
}
