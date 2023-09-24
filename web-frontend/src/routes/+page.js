import { redirect } from "@sveltejs/kit";
import { validateUser } from "../api";

export const ssr = false;
export const prerender = true;

export async function load({ fetch }) {
	let res;
	try {
		res = await validateUser(fetch);
	} catch (e) {
		return;
	}
	if (res.ok) {
		throw redirect(302, "/app");
	}
}
