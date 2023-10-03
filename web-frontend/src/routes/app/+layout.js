import { redirect } from "@sveltejs/kit";
import { loginUrl, validateUser } from "$lib/api";

export const prerender = true;

export async function load({ fetch, url }) {
	try {
		let res = await validateUser(fetch, url.origin);
		if (res.ok) {
			return {};
		}
	} catch (e) {}
	throw redirect(302, loginUrl(url.origin));
}
