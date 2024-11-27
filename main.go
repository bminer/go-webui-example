package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed web
var webFS embed.FS

func main() {
	webFSSub, err := fs.Sub(webFS, "web")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(webFSSub)))
	mux.HandleFunc("/webui.js", func(w http.ResponseWriter, r *http.Request) {
		log.Println("How to serve this?")
		http.NotFound(w, r)
	})

	// run the app on the specified port
	// Use the build tags to determine if this runs in "desktop" app mode or
	// "headless" server mode.
	run(8080, mux)
}
