package main

import (
	"filter/src/admission/server"
)

func main() {
	// Starts the server listening for incoming requests
	server.ListenAndServe()
}
