package config

import "os"

var (
	LabelValue = "sidecar"
	LabelKey   = os.Getenv("LABEL")

	TLS_CRT = os.Getenv("TLS_CRT")
	TLS_KEY = os.Getenv("TLS_KEY")
)
