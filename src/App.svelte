<script lang="ts">
	let video: HTMLVideoElement;
	let noSeek = false;
	let paused = true;

	let nono: Set<string> = new Set();

	interface Data {
		key: string;
	}

	interface SeekData {
		key: string;
		time: number;
	}

	const ws = new WebSocket("ws://127.0.0.1:8080/socket/1");
	ws.onopen = () => {
		console.log("connected");
	};
	ws.onerror = (e) => {
		console.log("error", e);
	};
	ws.onmessage = (event: MessageEvent<string>) => {
		console.log("got message:", event.data);
		const data = JSON.parse(event.data) as Data;
		switch (data.key) {
			case "PAUSE":
				nono.add("PAUSE");
				video.pause();
				break;
			case "PLAY":
				nono.add("PLAY");
				video.play();
				break;
			case "SEEKED":
				nono.add("PAUSE");
				video.pause();
				nono.add("SEEKED");
				video.currentTime = (data as SeekData).time;
				break;
		}
	};

	function send(data: Data | SeekData, all = false) {
		ws.send((all ? "!" : "") + JSON.stringify(data));
	}

	function onPause() {
		paused = true;

		if (nono.has("PAUSE")) {
			nono.delete("PAUSE");
			return;
		}
		console.log("on Pause");
		send({ key: "PAUSE" });
	}

	function onPlay(event) {
		paused = false;

		if (nono.has("PLAY")) {
			nono.delete("PLAY");
			return;
		}
		console.log("on Play");
		send({ key: "PLAY" });
	}

	function onSeeked(event) {
		if (nono.has("SEEKED")) {
			nono.delete("SEEKED");
			return;
		}

		console.log("on Seeked");
		send({ key: "SEEKED", time: video.currentTime });
		if (!paused) {
			send({ key: "PLAY" }, true);
		}
	}
</script>

<main>
	<!-- create video with caption track -->
	<video
		bind:this={video}
		on:pause={onPause}
		on:play={onPlay}
		on:seeked={onSeeked}
		id="video"
		width="640"
		height="360"
		controls
		autoplay
	>
		<source src="/assets/video.mp4" type="video/mp4" />
		<track kind="captions" label="English" src="" />
	</video>
</main>
