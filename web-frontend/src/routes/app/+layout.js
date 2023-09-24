import { redirect } from '@sveltejs/kit';

export async function load({ fetch, url }) {
	let res = await fetch('http://localhost:8000/auth/validate', {
		credentials: 'include'
	});
	if (res.ok) {
		return;
	}
	throw redirect(302, 'http://localhost:8000/auth/login?redirect_url=' + url.href);
}
