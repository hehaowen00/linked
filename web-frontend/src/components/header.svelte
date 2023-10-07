<script>
	import { logoutUrl } from "$lib/api";
	import { colorScheme } from "$lib/stores";
	export let url;

	let theme;

	function setTheme(value) {
		if (value === "dark") {
			$colorScheme = value;
		} else if (value === "light") {
			$colorScheme = value;
		}
	}

	colorScheme.subscribe((v) => {
		theme = v;
	});

	$: setTheme(theme);
</script>

<div class="row">
	<div class="row">
		<a href="/">Linked</a>
	</div>
	<div class="collapse-btn ml-auto">
		<a href="#">Menu</a>
	</div>
	<div class="row spaced-left ml-1">
		<a href="/app/">Bookmarks</a>
		<a href="/app/bookmarks/c9s">Collections</a>
	</div>
	<div class="row ml-auto spaced-left">
		<select on:change={setTheme} bind:value={theme}>
			<option value="dark">Dark Mode</option>
			<option value="light">Light Mode</option>
		</select>
		<a href={logoutUrl(url)}>Logout</a>
	</div>
</div>
