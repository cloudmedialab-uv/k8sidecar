package gen

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"time"
)

// GenCert generates an ECDSA certificate for a given name and returns its PEM-encoded representation
// alongside the PEM-encoded private key.
func GenCert(name string) (string, string, error) {
	// Generate a DNS name based on the provided name
	dnsName := "service-sidecar-" + name + ".default.svc"

	// Generate ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Printf("Error generating ECDSA private key: %v", err)
		return "", "", err
	}

	// Create a template for our certificate
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "Andoni Salcedo Navarro",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // 1 year
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{dnsName},
	}

	// Create the certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Printf("Error creating the certificate: %v", err)
		return "", "", err
	}

	// Create byte buffers to store the PEM-encoded key and certificate
	certBuffer := new(bytes.Buffer)
	keyBuffer := new(bytes.Buffer)

	// Save the certificate in PEM format
	if err := pem.Encode(certBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Printf("Error encoding the certificate in PEM format: %v", err)
		return "", "", err
	}

	// Save the private key in PEM format
	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		log.Printf("Error marshaling the private key: %v", err)
		return "", "", err
	}
	if err := pem.Encode(keyBuffer, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		log.Printf("Error encoding the private key in PEM format: %v", err)
		return "", "", err
	}

	return certBuffer.String(), keyBuffer.String(), nil
}
