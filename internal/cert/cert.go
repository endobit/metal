// Package cert is for creating signed and self-signed certificates.
package cert

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

const (
	ED25519 = "Ed25519"
	P224    = "P224"
	P256    = "P256"
	P384    = "P384"
	P521    = "P521"
	RSA     = "RSA"
)

func decodeCert(raw []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, fmt.Errorf("failed to parse certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	return cert, nil
}

func decodeKey(raw []byte) (any, error) {
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, fmt.Errorf("failed to parse key PEM")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return key, nil
}
