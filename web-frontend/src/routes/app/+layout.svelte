<script>
	import Header from "../../components/header.svelte";
</script>

<div class="main">
	<Header url={window.origin} />
	<slot />
</div>
