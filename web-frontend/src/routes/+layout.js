import { validateUser } from "$lib/api";

export const prerender = false;

export async function load({ fetch, url }) {
	let res;
	try {
		res = await validateUser(fetch, url.origin);
	} catch (e) {}
}
