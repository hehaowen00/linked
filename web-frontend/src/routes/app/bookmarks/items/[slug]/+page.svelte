<script>
	import { getOpenGraphInfo, putItem } from "$lib/api.js";
	import { displayTimestamp } from "$lib/util.js";

	export let data;
	let { item } = data;
	let { id, url, title, desc, created_at, deleted_at } = item;
	let isEditing = false;
	let titleValue = title;
	let loading = false;

	async function refresh() {
		loading = true;

		let res = await getOpenGraphInfo(window.fetch, window.origin, url);

		if (!res.ok) {
			loading = false;
			return;
		}

		let json = await res.json();
		titleValue = json.title;
		desc = json.desc;

		loading = false;
	}

	async function updateItem() {
		let res = await putItem(window.fetch, window.origin, {
			id,
			url,
			title: titleValue,
			desc,
			created_at,
			deleted_at
		});
		if (res.ok) {
			title = titleValue;
		}
		cancelEdit();
	}

	function editItem() {
		isEditing = true;
	}

	function cancelEdit() {
		titleValue = title;
		isEditing = false;
	}

	function cancel() {
		history.back();
	}

	function handleTextArea(e) {
		if (e.key === "Enter") {
			e.preventDefault();
		}
	}
</script>

<br />
<h3>{title}</h3>
<p />

{#if !isEditing}
	<div class="content">
		<div class="row">
			<a class="link" href={url} target="_blank">{url}</a>
		</div>
		<p />
		<div class="row">
			<span class="timestamp">Added {displayTimestamp(created_at)}</span>
			<p />
		</div>
		{#if desc}
			<div class="row">
				{desc}
			</div>
			<br />
		{/if}
		<div class="row">
			<span class="font-sm">Add To Collection</span>
		</div>
		<div class="row">
			<select class="w-100">
				<option>Collection 1</option>
			</select>
		</div>
		<br />
		<div class="row font-sm spaced-left">
			<span>Tags: </span>
			<span>#golang</span>
			<span>#go</span>
			<span>#programming</span>
		</div>
		<div class="row">
			<select class="w-100">
				<option>Tag 1</option>
			</select>
		</div>
		<br />
		<div class="row spaced-left">
			<button on:click={editItem}>Edit</button>
			<button on:click={cancel}>Cancel</button>
		</div>
	</div>
{:else}
	<div class="content">
		<div class="row">
			<a class="link" href={url} target="_blank">{url}</a>
		</div>
		<p />
		<!-- <div class="row">Collection</div> -->
		<!-- <div class="row"> -->
		<!-- 	<select class="w-100"> -->
		<!-- 		<option>Collection 1</option> -->
		<!-- 		<option>Collection 1</option> -->
		<!-- 		<option>Collection 1</option> -->
		<!-- 		<option>Collection 1</option> -->
		<!-- 		<option>Collection 1</option> -->
		<!-- 	</select> -->
		<!-- </div> -->
		<div class="row font-semibold">Title</div>
		<div class="row">
			<textarea
				rows={3}
				placeholder="Title"
				bind:value={titleValue}
				on:keydown={handleTextArea}
				spellcheck="false" />
		</div>
		<p />
		<div class="row font-semibold">
			<span>Description</span>
		</div>
		<div class="row">
			<textarea
				rows={4}
				placeholder="Description"
				bind:value={desc}
				on:keydown={handleTextArea}
				spellcheck="false" />
		</div>
		<br />
		<div class="row spaced-left">
			<button on:click={refresh} disabled={loading}>Refresh</button>
			<button on:click={updateItem} disabled={loading}>Save</button>
			<button on:click={cancelEdit} disabled={loading}>Cancel</button>
		</div>
	</div>
{/if}
