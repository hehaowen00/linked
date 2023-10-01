<script>
	import Header from "../../components/header.svelte";
	export let data;
</script>

<div class="main">
	<Header url={data.url.origin} />
	<slot />
</div>
