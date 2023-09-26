<script>
	import Header from "../../../../components/header.svelte";
	import { getOpenGraphInfo, postItem } from "../../../../api.js";
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

	async function pasteUrl() {
		let resp = await navigator.clipboard.readText();
		url = resp;
	}

	async function fetchOpenGraph(url) {
		disableAddButton = false;
		if (!checkValidUrl(url)) {
			og = {
				title: "",
				desc: "",
				image_url: ""
			};
			return;
		}

		try {
			let res = await getOpenGraphInfo(window.fetch, data.url.origin, url);
			if (res.ok) {
				og = await res.json();
			} else {
				console.log("error", await res.json());
			}
		} catch (e) {
			console.log(e);
		}
	}

	$: url && fetchOpenGraph(url);

	async function addItem() {
		if (!checkValidUrl(url) || og.title == "") {
			return;
		}
		let res = await postItem(window.fetch, data.url.origin, collection.id, {
			url,
			title: og.title,
			desc: og.desc
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
	<button on:click={addItem} disabled={og.title === "" || !checkValidUrl(url)}>Add Item</button>
</div>
<p />
{#if url}
	<div class="flex wrap">
		<a href={url}>
			{url}
		</a>
	</div>
	<p />
{/if}
<div class="flex flex-row">
	<input type="text" placeholder="Title" bind:value={og.title} />
</div>
<p />
<div class="flex flex-row">
	<input type="text" placeholder="Description" bind:value={og.desc} />
</div>
{#each items as item}
	<p>
		<a href={item.url} target="_blank">
			{item.title}
		</a>
	</p>
	{#if item.desc}
		<p class="flex flex-row item-desc">{item.desc}</p>
	{/if}
	<p>Added {displayTimestamp(item.created_at)}</p>
	<button>Edit</button>
	<button>Archive</button>
{/each}
