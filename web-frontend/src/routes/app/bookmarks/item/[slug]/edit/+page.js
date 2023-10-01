import { getItemById } from "../../../../../../api";

export const ssr = false;
export const prerender = false;

export async function load({ fetch, url, params }) {
	let res = await getItemById(fetch, url.origin, params.slug);
	if (!res.ok) {
		return {};
	}

	let json = await res.json();

	return {
		item: json.data,
		url
	};
}
