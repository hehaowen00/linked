import { getCollections } from "../../api";

export async function load({ fetch }) {
	let res = await getCollections(fetch);
	if (!res.ok) {
	}
	let json = await res.json();
	return {
		collections: json.data ?? []
	};
}
