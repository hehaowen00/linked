<script>
	import { getOpenGraphInfo, putItem } from "../../../../../../api.js";

	export let data;
	let { item } = data;
	let { id, url, title, desc, created_at, deleted_at } = item;
	let titleValue = title;
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
			titleValue = json.title;
			desc = json.desc;
		} catch (e) {
			loading = false;
		}
		loading = false;
	}

	async function updateItem() {
		let res = await putItem(window.fetch, data.url.origin, {
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
	}

	function back() {
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

<div class="content">
	<div class="row">
		<a href={url} target="_blank">{url}</a>
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
	<div class="row">Title</div>
	<div class="row">
		<textarea
			rows={3}
			placeholder="Title"
			bind:value={titleValue}
			on:keydown={handleTextArea}
			spellcheck="false" />
	</div>
	<p />
	<div class="row">
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
		<button on:click={back} disabled={loading}>Back</button>
	</div>
</div>
