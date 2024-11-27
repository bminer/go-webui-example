//go:build desktop

package main

import (
	"fmt"
	"net/http"

	ui "github.com/webui-dev/go-webui/v2"

	// This directive specifies linker flags for cgo.
	// -Wl,-Bstatic: Link with static libraries.
	// -lpthread: Link with the pthread library.
	// -Wl,-Bdynamic: Link with dynamic libraries.
	// -Wl,-subsystem,windows: Specify the Windows subsystem for the executable.

	// #cgo LDFLAGS: -Wl,-Bstatic -lpthread -Wl,-Bdynamic -Wl,-subsystem,windows
	"C"
)

func greet(e ui.Event) string {
	name, _ := ui.GetArg[string](e)
	fmt.Printf("%s has reached the backend!\n", name)
	jsResp := fmt.Sprintf("Hello %s 🐇", name)
	return jsResp
}

// run opens a web desktop application pointed to the HTTP server
func run(port int, h http.Handler) {
	portStr := fmt.Sprintf(":%d", port)

	// Launch HTTP server
	fmt.Println("Server is running on", portStr)
	go http.ListenAndServe(portStr, h)

	// Create a window
	w := ui.NewWindow()

	// Set port for window if we used bindings
	w.SetPort(8081)
	// Bind a Go function
	ui.Bind(w, "greet", greet)

	// Show frontend
	w.Show("http://localhost" + portStr)
	// Wait until all windows get closed.
	ui.Wait()
}
