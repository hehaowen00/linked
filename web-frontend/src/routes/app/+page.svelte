<script>
	import { onMount } from "svelte";

	export let data;
	let name = "";
	let collections = data.collections;

	async function fetchCollections() {
		let res = await fetch("http://localhost:8000/api/collections", {
			credentials: "include"
		});
		let json = await res.json();
		collections = json.data;
	}

	async function newCollection() {
		let res = await fetch("http://localhost:8000/api/collections", {
			method: "POST",
			credentials: "include",
			body: JSON.stringify({
				name: name
			})
		});
		if (!res.ok) {
			return;
		}
		name = "";
		await fetchCollections();
	}

	async function deleteCollection(index) {
		let c = collections[index];
		let res = await fetch(`http://localhost:8000/api/collections/${c.id}`, {
			method: "DELETE",
			credentials: "include",
			body: JSON.stringify({
				name: c.name
			})
		});
		if (!res.ok) {
			if (res.status === 401) {
			}
			return;
		}
		await fetchCollections();
	}

	onMount(async () => {});
</script>

<a href="http://localhost:8000/auth/logout"><p>Logout</p></a>
<h1>Collections</h1>

<input type="text" placeholder="New Collection" bind:value={name} />
<button on:click={newCollection}> Add Collection </button>

{#each collections as collection, idx}
	{#if collection.deleted_at == 0}
		<p>
			<a href="/app/collections/{collection.id}">{collection.name}</a>
			{collection.deleted_at}
			<button on:click={() => deleteCollection(idx)}>Delete</button>
		</p>
	{/if}
{/each}
<hr />
{#each collections as collection, idx}
	{#if collection.deleted_at != 0}
		<p>
			<a href="/app/collections/{collection.id}">{collection.name}</a>
			{collection.deleted_at}
		</p>
	{/if}
{/each}
