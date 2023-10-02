import { browser } from "$app/environment";
import { writable } from "svelte/store";

export const dialogStore = writable({
	name: "",
	cb: function () {}
});

export const colorScheme = writable(localStorage.getItem("theme") ?? "light");

colorScheme.subscribe((theme) => {
	if (browser) {
		localStorage.setItem("theme", theme);
	}
});

export const userStore = writable({});
