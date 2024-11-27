//go:build !desktop

package main

import (
	"fmt"
	"net/http"
)

// run simply runs the HTTP server
func run(port int, h http.Handler) {
	portStr := fmt.Sprintf(":%d", port)
	fmt.Println("Server is running on", portStr)
	http.ListenAndServe(portStr, h)
}
