export async function load({ fetch, params }) {
	let collection = {};
	let slug = params.slug;

	let res = await fetch(`http://localhost:8000/api/collections/${slug}`, {
		credentials: "include"
	});

	if (!res.ok) {
	}

	let resp = await res.json();
	collection = resp.data;

	res = await fetch(`http://localhost:8000/api/collections/${slug}/items`, {
		credentials: "include"
	});

	let items = await res.json();

	return {
		slug,
		collection,
		items
	};
}
