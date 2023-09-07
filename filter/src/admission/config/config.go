package config

import "os"

// Setting up configuration values. Used across multiple parts of the application.
var (
	// LabelValue defines the value for the sidecar label.
	LabelValue = "sidecar"

	// LabelKey retrieves the key for the sidecar label from environment variables.
	LabelKey = os.Getenv("LABEL")

	// TLS_CRT retrieves the TLS certificate from environment variables, likely for secure communications.
	TLS_CRT = os.Getenv("TLS_CRT")

	// TLS_KEY retrieves the TLS key from environment variables, paired with the certificate for secure communications.
	TLS_KEY = os.Getenv("TLS_KEY")
)
