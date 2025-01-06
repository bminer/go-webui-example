//go:build !headless

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	ui "github.com/webui-dev/go-webui/v2"
)

// Statically link the pthread library
// -Wl,-Bstatic: Statically link the following libraries
// -lpthread: Link (statically) with the pthread library
// -Wl,-Bdynamic: Dynamically link the following libraries

// #cgo LDFLAGS: -Wl,-Bstatic -lpthread -Wl,-Bdynamic
import "C"

// run opens a web desktop application pointed to the HTTP server
func run(mux *http.ServeMux) {
	log.Println("Running in WebUI mode")

	// Listen on next available TCP port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	defer listener.Close()

	// Create a window
	defer ui.Clean()
	ui.SetTimeout(5)
	w := ui.NewWindow()
	// w.SetSize(1024, 768)

	// Binding to empty string means we capture all events on all elements
	w.Bind("", eventHandler)
	// Redirect /webui.js requests to WebUI
	mux.HandleFunc("/webui.js", func(rw http.ResponseWriter, r *http.Request) {
		webUIPort, err := w.GetPort()
		if err != nil || webUIPort == 0 {
			http.Error(rw, "WebUI port unknown", http.StatusServiceUnavailable)
			return
		}
		rw.Header().Set(
			"Location", fmt.Sprintf("http://localhost:%d/webui.js", webUIPort),
		)
		rw.WriteHeader(http.StatusFound)
	})

	// Start HTTP server
	go func() {
		log.Fatal(http.Serve(listener, mux))
	}()
	log.Printf("HTTP server is running on http://localhost:%d", port)

	// Show WebUI, trying a few different preferred browsers
	url := fmt.Sprintf("http://localhost:%d/?webui=true", port)
	log.Println("Trying ChromiumBased browser")
	err = w.ShowBrowser(url, ui.ChromiumBased)
	if err != nil {
		log.Println("ChromiumBased failed, trying any available browser")
		err = w.Show(url)
	}
	if err != nil {
		w.Destroy()
		ui.Exit()
		ui.Clean()
		log.Fatal(err)
	}

	// Wait until all windows are closed
	webUIPort, err := w.GetPort()
	if err != nil {
		panic(err)
	}
	log.Printf("WebUI is running on http://localhost:%d", webUIPort)
	ui.Wait()
}

// eventHandler is called for all WebUI events
func eventHandler(e ui.Event) any {
	switch e.EventType {
	case ui.Disconnected:
		fmt.Printf("Disconnected event: %+v\n", e)
		ui.Exit()
	case ui.Connected:
		fmt.Printf("Connected event: %+v\n", e)
	case ui.MouseClick:
		if e.Element == "callWebUI" {
			e.Window.Run(
				`document.querySelector("#greeting").innerText =
				"Go says: Hello, " + document.querySelector("#name").value;`,
			)
		}
	case ui.Navigation:
		fmt.Printf("Navigation event: %+v\n", e)
		// Since we bind all events, following `href` links is blocked by WebUI.
		// To control the navigation, we need to use `Navigate()`.
		// e.Window.Navigate(target)
	case ui.Callback:
	}
	return nil
}
