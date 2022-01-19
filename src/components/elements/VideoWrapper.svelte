<script lang="ts">
    import { onMount } from "svelte";
    import SkipPrevious from "svelte-material-icons/SkipPrevious.svelte"
    import SkipNext from "svelte-material-icons/SkipNext.svelte"
    import PlayCircle from "svelte-material-icons/PlayCircle.svelte"
    import Pause from "svelte-material-icons/Pause.svelte"
    import Fullscreen from "svelte-material-icons/Fullscreen.svelte"
    import { username, ws } from "../../store";



	let video: HTMLVideoElement;
    let videoDuration;
    let videoProgress = 0;
    let videoWrapper: HTMLDivElement;
	let noSeek = false;
	let paused = true;
    let fullscreen = false;
    let bubble: HTMLOutputElement;

	let nono: Set<string> = new Set();

	interface Data {
		key: string;
	}

	interface SeekData {
		key: string;
		time: number;
	}

	$ws.onopen = () => {
		console.log("connection established..");

        $ws.send(`JOIN ${$username}`)
	};
	$ws.onerror = (e) => {
		console.log("error", e);
	};
	$ws.onmessage = (event) => {
		console.log("got message:", event.data);
		const data = JSON.parse(event.data) as Data;
		switch (data.key) {
			case "PAUSE":
		        paused = true;
				video.pause();
				break;
			case "PLAY":
		        paused = false;
				video.play();
				break;
			case "SEEKED":
				video.pause();
				video.currentTime = (data as SeekData).time;
				break;
		}
	};

	function send(data: Data | SeekData, all = false) {
		$ws.send((all ? "!" : "") + JSON.stringify(data));
	}

	function onPause() {
		paused = true;
        video.pause()

		console.log("on Pause");
		send({ key: "PAUSE" });
	}

	function onPlay(event) {
		paused = false;
        video.play()

		console.log("on Play");
		send({ key: "PLAY" });
	}

	function onSeeked(event) {

		console.log("on Seeked");
		send({ key: "SEEKED", time: video.currentTime });
		if (!paused) {
			send({ key: "PLAY" }, true);
		}
	}

    onMount(() => {


        // add an event listener for fullscreen changes
        videoWrapper.addEventListener('fullscreenchange', (e) => {

            if (document.fullscreenElement) {
                fullscreen = true;
            } else {
                fullscreen = false;
            }

        })
    })

    const toggleFullscreen = (e: MouseEvent) => {
        if (fullscreen) {
            document.exitFullscreen();

        } else {
            videoWrapper.requestFullscreen();
        }
    }

</script>

<div id="videoComponentWrapper">
    <div id="videoWrapper" bind:this={videoWrapper} >
        <video
            on:ended={e => paused = true}
            class:videoNotFullscreen={!fullscreen}
            bind:this={video}
            bind:duration={videoDuration}
            bind:currentTime={videoProgress}
            id="video"
        >
            <source src="/assets/video.mp4" type="video/mp4" />
            <track kind="captions" label="English" src="" />
        </video>
        <div id="videoControls">
            <div class="hoverPointer" on:click={e => console.log("SkipPrevious clicked")}>
                <SkipPrevious size="2rem" />
            </div>
            <div class="hoverPointer">
                {#if paused}
                <div on:click={onPlay} >
                    <PlayCircle size="2rem" />
                </div>
                {:else}
                <div on:click={onPause} >
                    <Pause size="2rem" />
                </div>
                {/if}
            </div>
            <div class="hoverPointer" on:click={e => console.log("SkipNext clicked")}>
                <SkipNext size="2rem" />
            </div>
            <div id="seekbarWrapper" class="hoverPointer" >
                <input type=range min="0" max={videoDuration} bind:value={videoProgress} on:change={onSeeked} step="0.01" />
                <output bind:this={bubble} class="bubble">{Math.round(videoProgress)}</output>
            </div>
            <div id="fullscreenIcon" class="hoverPointer" on:click={toggleFullscreen}>
                <Fullscreen size="2rem" />
            </div>
        </div>
    </div>
</div>

<style lang="scss">

    video {
        border-radius: .2rem;
    }
    .videoNotFullscreen {
        width: 55vw;
    }

    #videoWrapper {
        position: relative;

        transition: filter .2s;

        &:hover {
            transition: filter .2s;
            filter: drop-shadow(0rem 1rem 2rem black);

            #videoControls {
                visibility: visible;
            }
        }
    }

    #videoComponentWrapper {
        display: grid;
        place-content: center;
    }
    
    #videoControls {
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;

        visibility: hidden;
        display: flex;
        align-items: center;
        height: 3rem;

        background-color: #131313e8;
        column-gap: 1rem;
        padding: 0 1rem;

        #fullscreenIcon {
            margin-left: auto;
        }

    }

    .hoverPointer {
        &:hover {
            cursor: pointer;
        }
    }
    #seekbarWrapper {
        flex-grow: 1;
        margin: 0 2rem;

        display: flex;
        gap: 2rem;
        justify-content: center;

        input {
            width: 100%;
            margin: 0;
            padding: 0;
        }
    }

</style>
