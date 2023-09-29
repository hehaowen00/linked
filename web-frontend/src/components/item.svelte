<script>
	import { displayTimestamp } from "../util";

	export let canEdit = true;
	export let item;

	let { url, title, desc, created_at } = item;
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
</script>

<div class="item">
	<div class="row">
		<a href={url} target="_blank">
			{title}
		</a>
	</div>
	{#if desc}
		<p class="row item-desc">{desc}</p>
	{/if}
	<div class="timestamp">Added {displayTimestamp(created_at)}</div>
	<div class="row">
		<button on:click={copyLink}>
			{#if copied}
				Copied!
			{:else}
				Copy Link
			{/if}
		</button>

		{#if canEdit}
			<button>Edit</button>
			<button>Archive</button>
		{/if}
	</div>
	<br />
</div>
