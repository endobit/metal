package cert

// Code based on the command-line utility
// https://go.dev/src/crypto/tls/generate_cert.go and turn into a package.
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"time"
)

// Options contains the settings for creating a new certificate.
type Options struct {
	subject   string
	algorithm string
	hostname  string
	caCert    []byte
	caKey     []byte
	lifetime  time.Duration
}

// NewOptions returns an initialized Options.
func NewOptions(opts ...func(*Options)) *Options {
	options := Options{
		algorithm: RSA,
		lifetime:  time.Hour * 24 * 365,
	}

	for _, opt := range opts {
		opt(&options)
	}

	return &options
}

// WithLifetime is an option setting function for NewOptions. It sets the
// lifetime of the certificate. The default is 1 year.
func WithLifetime(d time.Duration) func(*Options) {
	return func(f *Options) {
		f.lifetime = d
	}
}

// WithCA is an option setting function for NewOptions. It sets tells the
// options to sign certificates with a certificate authority. The default is to
// issue self-signed certificates.
func WithCA(cert, key []byte) func(*Options) {
	return func(f *Options) {
		f.caCert = cert
		f.caKey = key
	}
}

// WithSubject is an option setting function for NewOptions. It sets the subject
// to s.
func WithSubject(s string) func(*Options) {
	return func(f *Options) {
		f.subject = s
	}
}

// Alogorithm is an option setting function for NewOptions. It sets the
// encryption algorithm, the default RSA.
func Alogorithm(c string) func(*Options) {
	return func(f *Options) {
		f.algorithm = c
	}
}

// Create creates a new self signed certificate.
func (o *Options) Create(cert, key io.Writer) error {
	var (
		priv        any
		err         error
		isCA, hasCA bool
	)

	if o.caCert != nil {
		hasCA = true
		if o.caKey == nil {
			return errors.New("missing CA private key")
		}
	}
	if o.caKey != nil {
		if o.caCert == nil {
			return errors.New("missing CA certificate")
		}
	}

	if !hasCA {
		isCA = true // going to sign myself
	}

	hostname := o.hostname
	if hostname == "" {
		hostname, err = os.Hostname()
		if err != nil {
			return fmt.Errorf("failed to get hostname: %w", err)
		}
	}

	switch o.algorithm {
	case ED25519:
		_, priv, err = ed25519.GenerateKey(rand.Reader)
	case P224:
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case P256:
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case P384:
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case P521:
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	case RSA:
		priv, err = rsa.GenerateKey(rand.Reader, 2048)
	default:
		return fmt.Errorf("unknown algorithm: %s", o.algorithm)
	}

	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	// ECDSA, ED25519 and RSA subject keys should have the DigitalSignature
	// KeyUsage bits set in the x509.Certificate template
	keyUsage := x509.KeyUsageDigitalSignature

	// Only RSA subject keys should have the KeyEncipherment KeyUsage bits set.
	// In the context of TLS this KeyUsage is particular to RSA key exchange and
	// authentication.
	if _, isRSA := priv.(*rsa.PrivateKey); isRSA {
		keyUsage |= x509.KeyUsageKeyEncipherment
	}

	if isCA {
		keyUsage |= x509.KeyUsageCertSign
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	now := time.Now()
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{o.subject},
			CommonName:   hostname,
		},
		NotBefore: now,
		NotAfter:  now.Add(o.lifetime),
		KeyUsage:  keyUsage,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		BasicConstraintsValid: true,
		IsCA:                  isCA,
		MaxPathLenZero:        isCA,
		IPAddresses:           ips(),
	}

	var (
		caCert *x509.Certificate
		caKey  any
	)

	if hasCA {
		caCert, err = decodeCert(o.caCert)
		if err != nil {
			return fmt.Errorf("failed to decode CA certificate: %w", err)
		}
		caKey, err = decodeKey(o.caKey)
		if err != nil {
			return fmt.Errorf("failed to decode CA key: %w", err)
		}
	} else {
		caCert = &template
		caKey = priv
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, caCert, publicKey(caKey), caKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	if err := pem.Encode(cert, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		return fmt.Errorf("failed to write data to cert.pem: %w", err)
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return fmt.Errorf("unable to marshal private key: %w", err)
	}
	if err := pem.Encode(key, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		return fmt.Errorf("failed to write data to key.pem: %w", err)
	}

	return nil
}

func ips() []net.IP {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil
	}

	ips := make([]net.IP, 0, len(addrs))
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			continue // skip invalid addresses
		}

		ips = append(ips, ip)
	}

	return ips
}

func publicKey(priv any) any {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey) //nolint:forcetypeassert
	default:
		return nil
	}
}
