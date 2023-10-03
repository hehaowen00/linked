<script>
	import { goto } from "$app/navigation";
	import { putCollection } from "$lib/api.js";

	export let data;

	let { collection, url } = data;
	let { id, name, created_at, deleted_at } = collection;

	async function saveCollection() {
		let res = await putCollection(window.fetch, url.origin, "PUT", {
			id,
			name,
			created_at,
			deleted_at
		});
	}

	function cancel() {
		goto("/app");
	}
</script>

<h1>Edit Collection</h1>
<p />

<div class="content">
	<div class="row">Name</div>
	<div class="row">
		<input type="text" placeholder="Name" bind:value={name} />
	</div>
	<p />
	<div class="row spaced-left">
		<button on:click={saveCollection}> Save </button>
		<button on:click={cancel}> Cancel </button>
	</div>
</div>
