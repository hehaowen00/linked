<script>
	let info = {
		email: "",
		password: ""
	};

	let error = "";

	async function login() {
		let res = await fetch("/auth/login", {
			method: "POST",
			body: JSON.stringify(info)
		});
		if (!res.ok) {
			let json = await res.json();
			error = json.error;
		}
	}
</script>

<div class="main col col-md text-center">
	<h1>Login</h1>
	<div class="col spaced-top-2">
		<input type="email" placeholder="Email Address" bind:value={info.email} />
		<input type="password" placeholder="Password" bind:value={info.password} />
	</div>
	<p class="text-right">
		<a href="/register">Don't have an account?</a>
	</p>
	{#if error}
		<p>{error}</p>
	{/if}
	<button on:click={login}>Login</button>
</div>
