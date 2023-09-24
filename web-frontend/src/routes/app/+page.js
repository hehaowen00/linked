export async function load({ fetch }) {
	let res = await fetch("http://localhost:8000/api/collections", {
		credentials: "include"
	});
	if (!res.ok) {
	}
	let json = await res.json();
	return {
		collections: json.data ?? []
	};
}
