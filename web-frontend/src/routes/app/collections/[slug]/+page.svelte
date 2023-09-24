<script>
	import { onMount } from 'svelte';

	export let data;
	let collection = data.collection ?? {};
	let items = [];
	let url = '';
	let og = {
		title: '',
		description: '',
		image_url: ''
	};

	async function fetchOpenGraph(url) {
		if (!url || (!url.startsWith('http://') && !url.startsWith('https://'))) {
			console.log('invalid url');
			og = {
				title: '',
				description: '',
				image_url: ''
			};
			return;
		}

		let res = await fetch('http://localhost:8000/api/opengraph/info', {
			method: 'POST',
			credentials: 'include',
			body: JSON.stringify({
				url: url
			})
		});

		og = await res.json();
	}

	$: url && fetchOpenGraph(url);

	async function addItem() {}

	onMount(async () => {
		let res = await fetch('http://localhost:8000/api/collections/' + collection.id + '/items', {
			credentials: 'include'
		});
		let resp = await res.json();
		if (resp.data) {
			items = resp.data;
		}
	});
</script>

<a href="/app"><p>Home</p></a>
<h1>{collection.name}</h1>
<input type="text" placeholder="URL" bind:value={url} />
<button on:click={addItem}>Add Item</button>
<a href={url}>
	<p>{url}</p>
</a>
<p>{og.title}</p>
<p>{og.description}</p>
