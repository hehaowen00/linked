import { redirect } from "@sveltejs/kit";

export const ssr = false;

export async function load({ fetch }) {
	let res;
	try {
		res = await fetch("http://localhost:8000/auth/validate", {
			credentials: "include"
		});
	} catch (e) {}
	console.log(res);
	if (res.ok) {
		throw redirect(302, "/app");
	}
}
