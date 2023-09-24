<script>
	import { onMount } from "svelte";
	import { getItems, getOpenGraphInfo, postItem } from "../../../../api.js";

	export let data;
	let collection = data.collection ?? {};
	let items = data.items ?? [];
	let url = "";
	let og = {
		title: "",
		description: "",
		image_url: ""
	};

	async function pasteUrl() {
		let resp = await navigator.clipboard.readText();
		url = resp;
	}

	async function fetchOpenGraph(url) {
		if (!url || (!url.startsWith("http://") && !url.startsWith("https://"))) {
			console.log("invalid url");
			og = {
				title: "",
				description: "",
				image_url: ""
			};
			return;
		}

		let res = await getOpenGraphInfo(window.fetch, url);
		og = await res.json();
	}

	$: url && fetchOpenGraph(url);

	async function addItem() {
		let res = await postItem(window.fetch, collection.id, {
			url,
			title: og.title,
			description: og.description
		});
		let resp = await res.json();
		items = [...items, resp.data];
	}

	onMount(async () => {
		let res = await getItems(window.fetch, collection.id);
		let resp = await res.json();
		if (resp.data) {
			items = resp.data;
		}
	});
</script>

<a href="/app"><p>Home</p></a>
<h1>{collection.name}</h1>
<input type="text" placeholder="URL" bind:value={url} />
<button on:click={pasteUrl}>Paste URL</button>
<button on:click={addItem}>Add Item</button>
<a href={url}>
	<p>{url}</p>
</a>
<p>
	<input type="text" placeholder="Title" bind:value={og.title} />
</p>
<p>
	<input type="text" placeholder="Description" bind:value={og.description} />
</p>

{#each items as item, i}
	<p>
		<a href={item.url}>
			{item.title}
		</a>
	</p>
	{#if item.description}
		<p>{item.description}</p>
	{/if}
	<p>Added {new Date(item.created_at).toString()}</p>
	<button>Delete</button>
{/each}
