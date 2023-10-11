<script>
	let info = {
		email: "",
		first: "",
		last: "",
		dob: undefined
	};
	let image = "";

	async function register() {
		let res = await fetch("/auth/register", {
			method: "POST",
			body: JSON.stringify(info)
		});
		if (!res.ok) {
			return;
		}

		let json = await res.json();
		image = json.qr_code;
	}

	let max = new Date().toISOString().slice(0, -14);
</script>

<div class="main col text-center col-md">
	<h1>Register</h1>
	{#if image}
	<div>
		<p>New User Created</p>
		<p>Scan QR Code with Authenticator App</p>
		<p>This is required to login</p>
		<img src={image} />
	</div>
	{:else}
	<div class="col spaced-top-2">
		<input type="email" placeholder="Email Address" bind:value={info.email} />
		<div class="row spaced-left">
			<input type="text" placeholder="First Name" bind:value={info.first} />
			<input type="text" placeholder="Last Name" bind:value={info.last} />
		</div>
		<input type="date" placeholder="Date of Birth" {max} bind:value={info.dob} />
	</div>
	<div>
		<p class="text-right">
			<a href="/login">Already have an account?</a>
		</p>
	</div>
	<button on:click={register}>Register</button>
	{/if}
</div>
