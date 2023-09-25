<script>
	export let collection;
	export let removeFn = function () {};

	let { id, name, created_at, deleted_at } = collection;

	function displayTimestamp(unixMillis) {
		return new Date(unixMillis).toLocaleString();
	}
</script>

<div>
	<p><a href="/app/collections/{id}">{name}</a></p>
	<p>Created At {displayTimestamp(created_at)}</p>
	{#if deleted_at == 0}
		<p><button on:click={() => removeFn()}>Delete</button></p>
	{:else}
		<p><button> Delete</button></p>
	{/if}
</div>
