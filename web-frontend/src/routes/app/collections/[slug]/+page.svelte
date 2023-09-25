<script>
	import { getOpenGraphInfo, postItem } from "../../../../api.js";
	import Header from "../../../../components/header.svelte";
	import { checkValidUrl, displayTimestamp } from "../../../../util.js";

	export let data;
	let collection = data.collection ?? {};
	let items = data.items ?? [];
	let url = "";
	let og = {
		title: "",
		desc: "",
		image_url: ""
	};
	let disabledAdd = false;

	async function pasteUrl() {
		let resp = await navigator.clipboard.readText();
		url = resp;
	}

	async function fetchOpenGraph(url) {
		if (!checkValidUrl(url)) {
			console.log("invalid url");
			og = {
				title: "",
				desc: "",
				image_url: ""
			};
			return;
		}

		disabledAdd = true;
		let res = await getOpenGraphInfo(window.fetch, data.url.origin, url);
		og = await res.json();
		disabledAdd = false;
	}

	$: url && fetchOpenGraph(url);

	async function addItem() {
		if (!checkValidUrl(url) || og.title == "") {
			return;
		}
		let res = await postItem(window.fetch, data.url.origin, collection.id, {
			url,
			title: og.title,
			desc: og.description
		});
		let resp = await res.json();
		items = [...items, resp.data];
		url = "";
		og = {
			title: "",
			desc: "",
			image_url: ""
		};
	}
</script>

<Header url={data.url.origin} />
<h1>{collection.name}</h1>
<div class="flex flex-row">
	<input type="text" placeholder="URL" bind:value={url} />
</div>
<p />
<div class="flex flex-row">
	<button on:click={pasteUrl}>Paste URL</button>
</div>
<p />
<div class="flex flex-row">
	<button on:click={addItem} disabled={disabledAdd}>Add Item</button>
</div>
<a href={url}>
	<p>{url}</p>
</a>
<div class="flex flex-row">
	<input type="text" placeholder="Title" bind:value={og.title} />
</div>
<p />
<div class="flex flex-row">
	<input type="text" placeholder="Description" bind:value={og.description} />
</div>
{#each items as item, i}
	<p>
		<a href={item.url}>
			{item.title}
		</a>
	</p>
	{#if item.desc}
		<p>{item.desc}</p>
	{/if}
	<p>Added {displayTimestamp(item.created_at)}</p>
	<button>Delete</button>
{/each}
