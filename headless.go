//go:build headless

package main

import (
	"log"
	"net"
	"net/http"
)

// run simply runs the HTTP server headlessly
func run(mux *http.ServeMux) {
	log.Println("Running in headless mode")

	// Listen on next available TCP port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	defer listener.Close()

	// Start HTTP server
	log.Printf("HTTP server is running on http://localhost:%d", port)
	log.Fatal(http.Serve(listener, mux))
}
