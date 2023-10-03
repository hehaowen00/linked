<script>
	import { loginUrl } from "$lib/api";
</script>

<div class="main col text-center">
	<h1>Linked</h1>
	<p>Bookmark Manager</p>
	<a href={loginUrl(window.origin, window.href)}> Sign In With Google </a>
</div>
