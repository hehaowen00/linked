import { getItems } from "$lib/api";

export const ssr = false;

export async function load({ fetch, url }) {
	let res = await getItems(fetch, url.origin);
	if (!res.ok) {
		return;
	}
	let json = await res.json();
	return {
		items: json.data ?? []
	};
}
