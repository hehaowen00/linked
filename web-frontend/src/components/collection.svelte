<script>
	export let collection;
	export let update;

	let { id, name, created_at, deleted_at } = collection;

	function displayTimestamp(unixMillis) {
		return new Date(unixMillis).toLocaleString();
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
</script>

<div>
	<p><a href="/app/collections/{id}">{name}</a></p>
	{#if deleted_at == 0}
		<p>Created At {displayTimestamp(created_at)}</p>
		<p>
			<button on:click={archive}>Edit</button>
			<button on:click={archive}>Archive</button>
		</p>
	{:else}
		<p>Archived At {displayTimestamp(deleted_at)}</p>
		<p>
			<button on:click={unarchive}>Unarchive</button>
			<button on:click={unarchive}>Delete</button>
		</p>
	{/if}
</div>
