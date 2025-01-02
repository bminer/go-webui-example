package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed web
var webFS embed.FS

func main() {
	production := !falsy(os.Getenv("PRODUCTION"))
	if production {
		log.Println("Environment: production")
	} else {
		log.Println("Environment: development")
	}

	webFSSub, err := fs.Sub(webFS, "web")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	// Expose environment variables
	mux.HandleFunc("/api/env", func(w http.ResponseWriter, r *http.Request) {
		// Set content type
		w.Header().Set("Content-Type", "application/json")
		// Write JSON response
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"production": production,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Expose public web files
	mux.Handle("/", http.FileServer(http.FS(webFSSub)))

	// Use the build tags to determine if this runs in "desktop" app mode or
	// "headless" server mode.
	run(mux)
}

// falsy returns true if the string is falsy
func falsy(s string) bool {
	return s == "" ||
		s == "0" ||
		strings.EqualFold(s, "false") ||
		strings.EqualFold(s, "null")
}
