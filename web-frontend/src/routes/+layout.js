import { validateUser } from "$lib/api";

export const prerender = false;

export async function load({ fetch, url }) {
	try {
		let res = await validateUser(fetch, url.origin);
	} catch (e) {}
}
