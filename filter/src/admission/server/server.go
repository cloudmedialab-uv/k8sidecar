package server

import (
	"crypto/tls"
	"filter/src/admission/config"
	"filter/src/admission/handlers"
	"log"
	"net/http"
)

func ListenAndServe() {
	// Load TLS certificate and key from configuration
	cert, err := tls.X509KeyPair([]byte(config.TLS_CRT), []byte(config.TLS_KEY))
	if err != nil {
		// Log the error if there's an issue loading the TLS cert and key, then exit
		log.Fatalf("Failed to load TLS certificate and key: %v", err)
	}

	// Define TLS configuration with loaded cert and minimum TLS version
	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	// Set up the HTTP server with given TLS configuration and bind it to port 8443
	server := &http.Server{
		Addr:      ":8443",
		Handler:   nil,
		TLSConfig: cfg,
	}

	log.Println("Setting up routes...")

	// Define the routes and their associated handlers
	http.HandleFunc("/kservice", handlers.KserviceHandler)
	http.HandleFunc("/deployment", handlers.DeploymentHandler)

	log.Println("Starting server on port 8443...")

	// Start the HTTPS server
	// Any error that arises from starting the server will be logged
	if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
		log.Printf("Error while starting server: %v\n", err)
		return
	}

	log.Println("Server stopped gracefully.")
}
