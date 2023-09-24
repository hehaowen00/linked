<script>
	import { onMount } from "svelte";
	import { getCollections, deleteCollection, postCollection, logoutUrl } from "../../api.js";

	export let data;
	let name = "";
	let collections = data.collections;

	async function newCollection() {
		let res = await postCollection(window.fetch, name);
		if (!res.ok) {
			return;
		}

		res = await getCollections(window.fetch);
		let json = await res.json();
		collections = json.data;
		name = "";
	}

	async function removeCollection(index) {
		let c = collections[index];
		let res = await deleteCollection(window.fetch, c);
		if (!res.ok) {
			if (res.status === 401) {
			}
			return;
		}
	}

	onMount(async () => {});
</script>

<a href={logoutUrl()}><p>Logout</p></a>
<h1>Collections</h1>

<input type="text" placeholder="New Collection" bind:value={name} />
<button on:click={newCollection}> Add Collection </button>

{#each collections as collection, idx}
	{#if collection.deleted_at == 0}
		<p>
			<a href="/app/collections/{collection.id}">{collection.name}</a>
			{collection.deleted_at}
			<button on:click={() => removeCollection(idx)}>Delete</button>
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
