<script>
	import DeleteDialog from "../components/deleteDialog.svelte";
	import { colorScheme } from "../stores";

	import "../styles/normalize.css";
	import "../styles/global.css";
	import "../styles/dark.css";
	import "../styles/light.css";

	let theme;

	function setTheme(value) {
		if (value === "dark" || value === "light") {
			document.documentElement.setAttribute("data-theme", value);
		}
	}

	colorScheme.subscribe((v) => {
		theme = v;
	});

	$: setTheme(theme);
</script>

<svelte:head>
	<title>Linked</title>
</svelte:head>

<div id="main-content">
	<DeleteDialog />
	<slot />
</div>
