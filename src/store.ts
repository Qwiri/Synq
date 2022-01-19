import { writable } from "svelte/store";
import type { Writable } from "svelte/store";

export const username: Writable<string> = writable('');
export const ws: Writable<WebSocket> = writable();