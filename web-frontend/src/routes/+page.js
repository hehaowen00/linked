import { redirect } from "@sveltejs/kit";
import { validateUser } from "$lib/api";

export const ssr = false;
export const prerender = true;

// export async function load({ fetch, url }) {
// 	let res;
// 	try {
// 		res = await validateUser(fetch, url.origin);
// 	} catch (e) {
// 		return {
// 			url
// 		};
// 	}
//
// 	if (res.ok) {
// 		throw redirect(302, "/app");
// 	}
// }
