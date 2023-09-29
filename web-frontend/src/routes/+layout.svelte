<script>
	import DeleteDialog from "../components/deleteDialog.svelte";

	import "../styles/normalize.css";
	import "../styles/global.css";
	import "../styles/dark.css";
	import "../styles/light.css";
	document.documentElement.setAttribute("data-theme", "light");
</script>

<div id="main-content">
	<DeleteDialog />
	<slot />
</div>
