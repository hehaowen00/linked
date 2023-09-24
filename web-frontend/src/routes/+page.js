import { redirect } from '@sveltejs/kit';

export const ssr = false;

export async function load({ fetch }) {
	let res = await fetch('http://localhost:8000/auth/validate', {
		credentials: 'include'
	});
	if (res.status == 200) {
		throw redirect(302, '/app');
	}
}
