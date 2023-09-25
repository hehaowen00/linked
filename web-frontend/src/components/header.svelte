<script>
	import { logoutUrl } from "../api";

	export let url;
	console.log("header", url);
</script>

<div class="row">
	<a href="/app">Home</a>
	<a class="right" href={logoutUrl(url)}>Logout</a>
</div>
