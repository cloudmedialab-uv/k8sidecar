package main

import (
	config "filter/admission/config"
	"filter/admission/server"
)

func main() {

	config := config.NewIntance()
	config.Load("LABEL_VALUE", "sidecar")
	config.Load("LABEL", "")
	config.Load("TLS_CRT", "")
	config.Load("TLS_KEY", "")
	config.Load("MEDIUM", "Memory")

	// Starts the server listening for incoming requests
	server.ListenAndServe()
}
