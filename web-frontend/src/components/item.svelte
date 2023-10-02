<script>
	import { goto } from "$app/navigation";
	import { removeItem } from "../api";
	import { displayTimestamp } from "../util";

	export let canEdit = true;
	export let item;
	export let collectionId;
	export let refresh;

	let { id, url, title, desc, created_at } = item;
	let copied = false;

	async function copyLink() {
		if (copied) {
			return;
		}
		await navigator.clipboard.writeText(url);
		copied = true;
		setTimeout(() => {
			copied = false;
		}, 500);
	}

	async function remove() {
		let res = await removeItem(window.fetch, window.origin, collectionId, id);
		if (!res.ok) {
			return;
		}
		await refresh();
	}

	function editItem() {
		goto(`/app/bookmarks/item/${id}/edit`);
	}
</script>

<div class="item">
	<div class="row">
		<a class="break-word" href={url} target="_blank">
			{title}
		</a>
	</div>
	<!-- {#if desc} -->
	<!-- 	<p class="row item-desc timestamp">{desc}</p> -->
	<!-- {/if} -->
	<div class="timestamp">Added {displayTimestamp(created_at)}</div>
	<div class="row spaced-left">
		<button on:click={copyLink}>
			{#if copied}
				Copied!
			{:else}
				Copy Link
			{/if}
		</button>

		{#if canEdit}
			<button on:click={editItem}>View</button>
			<button on:click={remove}>Remove</button>
		{/if}
	</div>
	<br />
</div>
