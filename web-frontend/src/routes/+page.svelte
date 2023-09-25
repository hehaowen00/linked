<script>
	import { loginUrl } from "../api.js";
	export let data;
	let { url } = data;
	console.log(url);
</script>

<div class="main">
	<h1>Linked</h1>
	<p>Bookmark Manager</p>
	<a href={loginUrl(url.origin, url.href)}> Sign In With Google </a>
</div>

<style>
	.main {
		width: 1000px;
		margin: 0 auto;
		text-align: center;
	}
</style>
