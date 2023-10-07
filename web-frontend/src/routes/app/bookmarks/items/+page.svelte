<script>
	import { getItems, getOpenGraphInfo, postItem } from "$lib/api.js";
	import { defaultOpenGraph } from "$lib/constants.js";
	import { checkValidUrl } from "$lib/util.js";

	let url = "";
	let opengraphInfo = { ...defaultOpenGraph };

	async function pasteUrl() {
		let resp = await navigator.clipboard.readText();
		if (!checkValidUrl(resp)) {
			return;
		}
		url = resp;
	}

	async function fetchOpenGraph() {
		if (!checkValidUrl(url)) {
			opengraphInfo = { ...defaultOpenGraph };
			return;
		}

		let check = new URL(url);
		if (!check.protocol.startsWith("http")) {
			return;
		}

		let res = await getOpenGraphInfo(window.fetch, window.origin, url);
		if (res.redirected) {
			goto(res.url);
		}
		if (res.ok) {
			opengraphInfo = await res.json();
		}
	}

	$: url && fetchOpenGraph(url);

	async function addItem() {
		if (!checkValidUrl(url) || opengraphInfo.title == "") {
			return;
		}
		let res = await postItem(window.fetch, window.origin, collection.id, {
			url,
			title: opengraphInfo.title,
			desc: opengraphInfo.desc
		});
		if (!res.ok) {
			return;
		}

		url = "";
		opengraphInfo = { ...defaultOpenGraph };

		await refresh();
	}

	async function refresh() {
		let res = await getItems(window.fetch, window.origin, collection.id);
		items = await res.json();
	}
</script>

<div class="content">
	<div class="row">
		<input type="text" placeholder="URL" bind:value={url} />
	</div>
	{#if url}
		<div class="row" style="font-size: 0.9rem;">
			<a class="link" href={url}>
				{url}
			</a>
		</div>
	{/if}
	<div class="row">
		<input type="text" placeholder="Title" bind:value={opengraphInfo.title} />
	</div>
	<div class="row">
		<input type="text" placeholder="Description" bind:value={opengraphInfo.desc} />
	</div>
	<div class="row">
		<div class="col">
			<button on:click={pasteUrl}>Paste URL</button>
		</div>
		<div class="col">
			<button on:click={addItem} disabled={opengraphInfo.title === "" || !checkValidUrl(url)}>
				Add Item
			</button>
		</div>
	</div>
</div>
<br />
