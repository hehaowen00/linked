import { PUBLIC_API_HOST } from "$env/static/public";

const API_HOST = PUBLIC_API_HOST;

export async function getCollectionById(fetch, id) {
	let res = await fetch(`http://${API_HOST}/api/collections/${id}`, {
		credentials: "include"
	});
	return res;
}

export async function getCollections(fetch) {
	let res = await fetch(`http://${API_HOST}/api/collections`, {
		credentials: "include"
	});
	return res;
}

export async function postCollection(fetch, name) {
	let res = await fetch(`http://${API_HOST}/api/collections`, {
		method: "POST",
		credentials: "include",
		body: JSON.stringify({
			name
		})
	});
	return res;
}

export async function deleteCollection(fetch, collection) {
	let res = await fetch(`http://${API_HOST}/api/collections/${collection.id}`, {
		method: "DELETE",
		credentials: "include",
		body: JSON.stringify({
			name: collection.name
		})
	});
	return res;
}

export async function getItems(fetch, id) {
	let res = await fetch(`http://${API_HOST}/api/collections/${id}/items`, {
		credentials: "include"
	});
	return res;
}

export async function postItem(fetch, collectionId, payload) {
	let res = await fetch(`http://${API_HOST}/api/collections/${collectionId}/items`, {
		method: "POST",
		credentials: "include",
		body: JSON.stringify(payload)
	});
	return res;
}

export async function getOpenGraphInfo(fetch, url) {
	let res = await fetch(`http://${API_HOST}/api/opengraph/info`, {
		method: "POST",
		credentials: "include",
		body: JSON.stringify({
			url
		})
	});
	return res;
}

export async function validateUser(fetch) {
	let res = await fetch(`http://${API_HOST}/auth/validate`, {
		credentials: "include"
	});
	return res;
}

export function loginUrl(redirect) {
	return `http://${API_HOST}/auth/login?redirect_url=${redirect}`;
}

export function logoutUrl() {
	return `http://${API_HOST}/auth/logout`;
}
