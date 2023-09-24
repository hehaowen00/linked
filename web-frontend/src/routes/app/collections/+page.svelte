<script>
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { onMount } from 'svelte';

	export let data;
	let collections = data.collections;

	async function newCollection(name) {
		let res = await fetch('http://localhost:8000/api/collections', {
			method: 'POST',
			credentials: 'include',
			body: JSON.stringify({
				name: name
			})
		});
		if (!res.ok) {
			return;
		}
		let resp = await res.json();
		collections.push(resp.data);
	}

	onMount(async () => {
		// let res = await fetch('http://localhost:8000/api/collections', {
		// 	credentials: 'include'
		// });
		// if (res.status === 401) {
		// 	console.log($page.url.href);
		// 	goto('http://localhost:8000/auth/login?redirect_url=' + $page.url.href);
		// } else {
		// 	console.log(await res.json());
		// }
	});
</script>

<h1>Collections</h1>

{#each collections as collection}
	<p>{collection.name}</p>
{/each}
