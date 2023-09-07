package server

import (
	"crypto/tls"
	"filter/src/admission/config"
	"filter/src/admission/handlers"
	"log"
	"net/http"
)

func ListenAndServe() {

	cert, err := tls.X509KeyPair([]byte(config.TLS_CRT), []byte(config.TLS_KEY))
	if err != nil {
		log.Fatalf("tls.X509KeyPair: %v", err)
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      ":8443",
		Handler:   nil,
		TLSConfig: cfg,
	}

	http.HandleFunc("/kservice", handlers.KserviceHandler)
	http.HandleFunc("/deployment", handlers.DeploymentHandler)

	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Println(err)
		return
	}

}
