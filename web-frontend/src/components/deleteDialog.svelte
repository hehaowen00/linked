<script>
	import { dialogStore } from "../stores";

	let dialog;
	let type = "";
	let name = "";
	let confirmValue = "";
	let cb = function () {
		console.log("callback");
	};

	dialogStore.subscribe((store) => {
		name = store.name;
		type = store.type;
		cb = store.cb;
		if (name) {
			dialog?.showModal();
			confirmValue = "";
		}
	});

	function cancel() {
		confirmValue = "";
		$dialogStore.name = "";
		$dialogStore.type = "";
		$dialogStore.cb = function () {};
	}
</script>

<dialog bind:this={dialog}>
	<div>
		Delete {type} "{name}" ?
	</div>
	<p />
	<div class="flex row">
		<input type="text" placeholder="Enter {name} to confirm" bind:value={confirmValue} />
	</div>
	<p />
	<form method="dialog">
		<div class="row">
			<div class="col">
				<button disabled={name !== confirmValue} on:click={cb}>Confirm</button>
			</div>
			<div class="col">
				<button on:click={cancel}>Cancel</button>
			</div>
		</div>
	</form>
</dialog>
