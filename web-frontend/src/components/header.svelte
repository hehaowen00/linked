<script>
	import { logoutUrl } from "../api";
	import { colorScheme } from "../stores";
	export let url;

	let theme;

	function setThemeDark() {
		document.documentElement.setAttribute("data-theme", "dark");
		$colorScheme = "dark";
	}

	function setThemeLight() {
		document.documentElement.setAttribute("data-theme", "light");
		$colorScheme = "light";
	}

	colorScheme.subscribe((v) => {
		theme = v;
	});

	$: setTheme(theme);

	function setTheme(value) {
		if (value === "dark") {
			setThemeDark();
		} else if (value === "light") {
			setThemeLight();
		}
	}
</script>

<div class="row">
	<div class="col">
		<a href="/app">Home</a>
	</div>
	<div class="row">
		<select on:change={setTheme} bind:value={theme}>
			<option value="dark">Dark Mode</option>
			<option value="light">Light Mode</option>
		</select>
		<a href={logoutUrl(url)}>Logout</a>
	</div>
</div>
