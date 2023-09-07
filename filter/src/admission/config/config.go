package config

import "os"

var (
	LabelKey   = os.Getenv("LABEL")
	LabelValue = "sidecar"

	TLS_CRT = os.Getenv("TLS_CRT")
	TLS_KEY = os.Getenv("TLS_KEY")
)
