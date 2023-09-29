<script>
	import { getCollections, putCollection, postCollection } from "../../api.js";
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
			console.log(await res.json());
			return;
		}

		await refresh();
		name = "";
	}

	async function update(method, collection) {
		let res = await putCollection(window.fetch, url.origin, method, collection);
		if (!res.ok) {
			console.log(await res.json());
			return;
		}
		await refresh();
	}

	async function refresh() {
		let res = await getCollections(window.fetch, url.origin);
		let json = await res.json();
		if (!res.ok) {
			console.log(await res.json());
			return;
		}
		collections = json.data;
	}
</script>

<Header url={url.origin} />

<h1>Collections</h1>
<p />

<div class="content">
	<div class="row">
		<input type="text" placeholder="New Collection" bind:value={name} />
	</div>
	<div class="row">
		<button on:click={newCollection}> Add Collection </button>
	</div>
	<p />
</div>

<div class="row">
	<select bind:value={category}>
		<option value="all" selected>All</option>
		<option value="deleted">Archived</option>
	</select>
</div>
<br />

{#if category == "all"}
	{#each collections as collection, idx}
		{#if collection.deleted_at == 0}
			<Collection bind:collection={collections[idx]} {update} />
		{/if}
	{/each}
{:else}
	{#each collections as collection, idx}
		{#if collection.deleted_at != 0}
			<Collection bind:collection={collections[idx]} {update} />
		{/if}
	{/each}
{/if}
