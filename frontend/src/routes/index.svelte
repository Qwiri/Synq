<script lang="ts">
	import Chat from "$lib/elements/Chat.svelte";
	import VideoWrapper from "$lib/elements/VideoWrapper.svelte";
	import Swal from 'sweetalert2/dist/sweetalert2.all.min.js';
	import { onMount } from "svelte";
	import { username, ws } from "../store";

	let connected = false;

	onMount(() => {

		// get username
		Swal.fire({
			title: 'Please enter a username:',
			icon: 'question',
			confirmButtonText: 'Cool',
			input: "text",
			background: "#181818",
			inputValidator: input => (!input && "You need to enter something")
		}).then(result => {
			$username = result.value;
			$ws = new WebSocket("ws://127.0.0.1:8080/socket/1");
			connected = true;

		})


	})

</script>

<main>
	{#if connected}
		<VideoWrapper />
		<Chat />
	{/if}
</main>

<style lang="scss">

	main {
		display: flex;
		width: 100vw;
		max-width: 100%;
		height: 100vh;

		* {
			flex-shrink: 0;
		}

		:global(#videoComponentWrapper) {
			width: 75vw;

		}
		:global(#chatWrapper) {
			width: 25vw;
		}
	}

</style>