//go:build !headless && production

package main

// Specify the Windows subsystem for the executable to avoid opening a console
// window.
// https://learn.microsoft.com/en-us/cpp/build/reference/subsystem-specify-subsystem?view=msvc-160

// #cgo LDFLAGS: -Wl,-subsystem,windows
import "C"
import "os"

// Forcefully set PRODUCTION environment variable
func init() {
	os.Setenv("PRODUCTION", "true")
}
