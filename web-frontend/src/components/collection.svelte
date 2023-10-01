<script>
	import { goto } from "$app/navigation";
	import { dialogStore } from "../stores";

	export let collection;
	export let update;

	let { id, name, created_at, deleted_at } = collection;

	function displayTimestamp(unixMillis) {
		return new Date(unixMillis).toLocaleString();
	}

	function gotoEdit() {
		goto(`/app/bookmarks/c8n/${id}/edit`);
	}

	async function unarchive() {
		collection.deleted_at = 0;
		let res = await update("PUT", collection);
		let json = await res.json();
		collection.deleted_at = json.data.deleted_at;
	}

	async function archive() {
		let res = await update("DELETE", collection);
		let json = await res.json();
		collection.deleted_at = json.data.deleted_at;
	}

	function confirm() {
		$dialogStore.type = "Collection";
		$dialogStore.name = name;
		$dialogStore.cb = async function () {
			let res = await update("DELETE", collection);
			let json = await res.json();
			collection.deleted_at = json.data.deleted_at;
		};
	}
</script>

<div class="collection">
	<div class="row">
		<a href="/app/bookmarks/c8n/{id}">{name}</a>
	</div>
	{#if deleted_at == 0}
		<div class="row spaced-left">
			<span class="timestamp">Created at {displayTimestamp(created_at)}</span>
		</div>
		<div class="row spaced-left">
			<button on:click={gotoEdit}>Edit</button>
			<button on:click={archive}>Archive</button>
		</div>
	{:else}
		<div class="timestamp">Archived at {displayTimestamp(deleted_at)}</div>
		<div class="row spaced-left">
			<button on:click={unarchive}>Unarchive</button>
			<button on:click={confirm}>Delete</button>
		</div>
	{/if}
	<br />
</div>
