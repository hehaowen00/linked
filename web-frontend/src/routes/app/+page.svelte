<script>
	import { getCollections, putCollection, postCollection } from "../../api.js";
	import Collection from "../../components/collection.svelte";

	export let data;
	let { collections, url } = data;
	let name = "";
	let category = "all";

	async function newCollection() {
		if (!name) {
			return;
		}

		let res = await postCollection(window.fetch, window.origin, name);
		if (!res.ok) {
			return;
		}

		await refresh();
		name = "";
	}

	async function update(method, collection) {
		let res = await putCollection(window.fetch, window.origin, method, collection);
		if (!res.ok) {
			return;
		}
		await refresh();
	}

	async function refresh() {
		let res = await getCollections(window.fetch, window.origin);
		let json = await res.json();
		if (!res.ok) {
			return;
		}
		collections = json.data;
	}
</script>

<h1>Bookmarks</h1>
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
	<select class="w-100" bind:value={category}>
		<option value="all" selected>All</option>
		<option value="deleted">Archived</option>
	</select>
</div>
<br />

{#if category == "all"}
	{#each collections as collection}
		{#key collection}
			<Collection {collection} {update} />
		{/key}
	{/each}
{:else}
	{#each collections as collection}
		{#if collection.deleted_at != 0}
			{#key collection}
				<Collection {collection} {update} />
			{/key}
		{/if}
	{/each}
{/if}
