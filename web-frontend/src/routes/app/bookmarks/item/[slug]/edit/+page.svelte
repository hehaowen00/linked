<script>
	import { getOpenGraphInfo, putItem } from "../../../../../../api.js";

	export let data;
	let { item } = data;
	let { id, url, title, desc, created_at, deleted_at } = item;
	let loading = false;

	async function refresh() {
		loading = true;
		let res = await getOpenGraphInfo(window.fetch, data.url.origin, url);
		if (!res.ok) {
			loading = false;
			return;
		}
		try {
			let json = await res.json();
			desc = json.desc;
		} catch (e) {
			loading = false;
		}
		loading = false;
	}

	async function updateItem() {
		await putItem(window.fetch, data.url.origin, {
			id,
			url,
			title,
			desc,
			created_at,
			deleted_at
		});
	}

	function back() {
		history.back();
	}
</script>

<h1>{title}</h1>
<p />

<div class="content">
	<div class="row">
		<a href={url} target="_blank">{url}</a>
	</div>
	<div class="row">
		<input type="text" placeholder="Description" bind:value={desc} />
	</div>
	<div class="row">
		<button on:click={refresh} disabled={loading}>Refresh</button>
		<button on:click={updateItem} disabled={loading}>Save</button>
		<button on:click={back} disabled={loading}>Back</button>
	</div>
	<div class="row">Move to collection</div>
	<div class="row">
		<select>
			<option>Collection 1</option>
			<option>Collection 1</option>
			<option>Collection 1</option>
			<option>Collection 1</option>
			<option>Collection 1</option>
		</select>
	</div>
</div>
