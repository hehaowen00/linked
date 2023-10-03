<script>
	import { goto } from "$app/navigation";
	import { displayTimestamp } from "$lib/util.js";

	export let data;
	let { item } = data;
	let { id, url, title, desc, created_at, deleted_at } = item;

	function editItem() {
		goto(`/app/bookmarks/item/${id}/edit`);
	}
</script>

<br />
<h3>{title}</h3>
<p />

<div class="content">
	<div class="row">
		<a class="link" href={url} target="_blank">{url}</a>
	</div>
	<p />
	<div class="row">
		<span class="timestamp">Added {displayTimestamp(created_at)}</span>
		<p />
	</div>
	<div class="row">
		{desc}
	</div>
	<br />
	<div class="row spaced-left">
		<button on:click={editItem}>Edit</button>
	</div>
</div>
