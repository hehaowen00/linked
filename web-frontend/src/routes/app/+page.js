import { redirect } from "@sveltejs/kit";
import { getCollections, loginUrl } from "../../api";

export const ssr = false;

export async function load({ fetch, url }) {
	try {
		let res = await validateUser(fetch, url.origin);
		if (!res.ok) {
			throw redirect(302, loginUrl(url.origin, url.href));
		}
	} catch (e) {}

	let res = await getCollections(fetch, url.origin);
	if (!res.ok) {
		return {
			url
		};
	}
	let json = await res.json();
	return {
		collections: json.data ?? [],
		url
	};
}
