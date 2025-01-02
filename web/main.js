const WEBUI_CONNECTION_TIMEOUT = 2 * 1000 // 2 seconds
const d = document
const w = window
d.addEventListener("DOMContentLoaded", async () => {
	const urlParams = new URLSearchParams(w.location.search)
	let env = null
	try {
		const res = await fetch("/api/env", {
			method: "GET",
			headers: {
				"Content-Type": "application/json",
			},
		})
		if (!res.ok) {
			throw new Error("status code " + res.status)
		}
		env = await res.json()
		console.log("Environment:", env)
	} catch (err) {
		console.error("Failed to fetch environment:", res.status, res.statusText)
	}

	if (urlParams.has("webui")) {
		function fatal(msg) {
			console.error(msg)
			w.close()
			alert("Fatal error: " + msg)
		}
		if (env == null) {
			fatal("Failed to fetch environment")
		}
		// Close window if connection times out
		const timer = setTimeout(
			() => fatal("Connection timed out"),
			WEBUI_CONNECTION_TIMEOUT,
		)
		// Setup WebUI event handlers
		webui.setEventCallback(async (e) => {
			if (e == webui.event.CONNECTED) {
				clearTimeout(timer)
				console.log("WebUI Connected!")
			} else if (e == webui.event.DISCONNECTED) {
				fatal("WebUI Disconnected!")
			}
		})
		if (env.production) {
			// Disable right-click and debugging tools in production
			d.addEventListener("contextmenu", (ev) => ev.preventDefault())
			d.addEventListener("keydown", (ev) => {
				if (
					ev.key.match(/^F\d+$/) ||
					(ev.ctrlKey && ev.shiftKey && ev.key === "I")
				) {
					ev.preventDefault()
				}
			})
		}
	}
})
