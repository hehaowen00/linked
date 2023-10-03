<script>
	import { page } from "$app/stores";
</script>

<div class="main">
	<h1>{$page.status} {$page.error.message}</h1>
	<a href="/">Click here to go home</a>
</div>
