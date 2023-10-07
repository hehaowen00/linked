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

export async function putCollection(fetch, host, method, collection) {
	let res = await fetch(`${host}/api/collections/${collection.id}`, {
		method,
		credentials: "include",
		body: JSON.stringify(collection)
	});
	return res;
}

export async function getItemsByCollection(fetch, host, id) {
	let res = await fetch(`${host}/api/collections/${id}/items`, {
		credentials: "include"
	});
	return res;
}

export async function getItems(fetch, host) {
	let res = await fetch(`${host}/api/items`, {
		credentials: "include"
	});
	return res;
}

export async function getItemById(fetch, host, id) {
	let res = await fetch(`${host}/api/items/${id}`, {
		credentials: "include"
	});
	return res;
}

export async function postItem(fetch, host, payload) {
	let res = await fetch(`${host}/api/items`, {
		method: "POST",
		credentials: "include",
		body: JSON.stringify(payload)
	});
	return res;
}

export async function putItem(fetch, host, item) {
	let res = await fetch(`${host}/api/items/${item.id}`, {
		method: "PUT",
		credentials: "include",
		body: JSON.stringify(item)
	});
	return res;
}

export async function removeItem(fetch, host, itemId) {
	let res = await fetch(`${host}/api/items/${itemId}`, {
		method: "DELETE",
		credentials: "include"
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
	let url = new URL(`${host}/auth/login`);
	if (redirect) {
		url.searchParams.append("redirect_url", redirect);
	}
	return url.toString();
}

export function logoutUrl(host) {
	return `${host}/auth/logout`;
}
