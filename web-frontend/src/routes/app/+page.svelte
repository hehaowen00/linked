<script>
	import { getCollections, deleteCollection, postCollection } from "../../api.js";
	import Collection from "../../components/collection.svelte";
	import Header from "../../components/header.svelte";

	export let data;
	let { collections, url } = data;
	let name = "";
	let category = "all";

	async function newCollection() {
		if (!name) {
			return;
		}

		let res = await postCollection(window.fetch, url.origin, name);
		if (!res.ok) {
			return;
		}

		res = await getCollections(window.fetch, url.origin);
		let json = await res.json();
		collections = json.data;
		name = "";
	}

	async function removeCollection(index) {
		let c = collections[index];
		let res = await deleteCollection(window.fetch, url.origin, c);
		if (!res.ok) {
			if (res.status === 401) {
			}
			return;
		}
	}
</script>

<Header url={url.origin} />

<h1>Collections</h1>

<div class="flex flex-row">
	<input type="text" placeholder="New Collection" bind:value={name} />
</div>
<p />
<div class="flex flex-row">
	<button on:click={newCollection}> Add Collection </button>
</div>
<p />

<div class="row">
	<select bind:value={category}>
		<option value="all" selected>All</option>
		<option value="active">Active</option>
		<option value="deleted">Deleted</option>
	</select>
</div>

{#if category == "all"}
	{#each collections as collection, idx}
		<Collection {collection} removeFn={() => removeCollection(idx)} />
	{/each}
{:else if category == "active"}
	{#each collections as collection, idx}
		{#if collection.deleted_at == 0}
			<Collection {collection} removeFn={() => removeCollection(idx)} />
		{/if}
	{/each}
{:else}
	{#each collections as collection, idx}
		{#if collection.deleted_at != 0}
			<Collection {collection} removeFn={() => removeCollection(idx)} />
		{/if}
	{/each}
{/if}

<style>
</style>
