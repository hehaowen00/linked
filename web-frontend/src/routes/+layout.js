import { validateUser } from "../api";

export const prerender = false;

export async function load({ fetch, url }) {
	let res;
	try {
		res = await validateUser(fetch, url.origin);
		let json = await res.json();

		return {
			info: json.data,
			url
		};
	} catch (e) {
		return {
			url
		};
	}
}
