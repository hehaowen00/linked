export async function getCollectionById(fetch, host, id) {
	let res = await fetch(`${host}/api/collections/${id}`, {
		credentials: "include"
	});
	return res;
}

export async function getCollections(fetch, host) {
	let res = await fetch(`${host}/api/collections`, {
		credentials: "include"
	});
	return res;
}

export async function postCollection(fetch, host, name) {
	let res = await fetch(`${host}/api/collections`, {
		method: "POST",
		credentials: "include",
		body: JSON.stringify({
			name
		})
	});
	return res;
}

export async function deleteCollection(fetch, host, collection) {
	let res = await fetch(`${host}/api/collections/${collection.id}`, {
		method: "DELETE",
		credentials: "include",
		body: JSON.stringify({
			name: collection.name
		})
	});
	return res;
}

export async function getItems(fetch, host, id) {
	let res = await fetch(`${host}/api/collections/${id}/items`, {
		credentials: "include"
	});
	return res;
}

export async function postItem(fetch, host, collectionId, payload) {
	let res = await fetch(`${host}/api/collections/${collectionId}/items`, {
		method: "POST",
		credentials: "include",
		body: JSON.stringify(payload)
	});
	return res;
}

export async function getOpenGraphInfo(fetch, host, url) {
	let res = await fetch(`${host}/api/opengraph/info`, {
		method: "POST",
		credentials: "include",
		body: JSON.stringify({
			url
		})
	});
	return res;
}

export async function validateUser(fetch, host) {
	let res = await fetch(`${host}/auth/validate`, {
		credentials: "include"
	});
	return res;
}

export function loginUrl(host, redirect) {
	return `${host}/auth/login?redirect_url=${redirect}`;
}

export function logoutUrl(host) {
	return `${host}/auth/logout`;
}
