// Copyright 2020-2025 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gufodao

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v2"
)

// GenerateCertificates creates self-signed CA, server and client certificates.
func GenerateCertificates(c *cli.Context) error {
	outputDir := "./certs"
	if c.Args().Len() > 0 {
		outputDir = c.Args().First()
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	caCertPath := filepath.Join(outputDir, "ca.pem")
	caKeyPath := filepath.Join(outputDir, "ca-key.pem")
	serverCertPath := filepath.Join(outputDir, "server.pem")
	serverKeyPath := filepath.Join(outputDir, "server-key.pem")
	clientCertPath := filepath.Join(outputDir, "client.pem")
	clientKeyPath := filepath.Join(outputDir, "client-key.pem")

	fmt.Printf("🔧 Generating certificates in %s ...\n", outputDir)

	// 1️⃣ Generate CA key and certificate
	caPriv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	caTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Gufo CA"},
			CommonName:   "Gufo Root CA",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caDER, _ := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caPriv.PublicKey, caPriv)
	writeCert(caCertPath, caDER)
	writeKey(caKeyPath, caPriv)

	// 2️⃣ Generate Server cert
	serverPriv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	serverTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject: pkix.Name{
			Organization: []string{"Gufo Server"},
			CommonName:   "gufo-server",
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(2, 0, 0),
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:        false,
	}
	serverDER, _ := x509.CreateCertificate(rand.Reader, serverTemplate, caTemplate, &serverPriv.PublicKey, caPriv)
	writeCert(serverCertPath, serverDER)
	writeKey(serverKeyPath, serverPriv)

	// 3️⃣ Generate Client cert
	clientPriv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	clientTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject: pkix.Name{
			Organization: []string{"Gufo Client"},
			CommonName:   "gufo-client",
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(2, 0, 0),
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		IsCA:        false,
	}
	clientDER, _ := x509.CreateCertificate(rand.Reader, clientTemplate, caTemplate, &clientPriv.PublicKey, caPriv)
	writeCert(clientCertPath, clientDER)
	writeKey(clientKeyPath, clientPriv)

	fmt.Println("✅ Certificates generated successfully!")
	fmt.Printf("CA: %s\nServer: %s\nClient: %s\n", caCertPath, serverCertPath, clientCertPath)
	return nil
}

// writeCert writes PEM-encoded certificate to file
func writeCert(path string, certDER []byte) {
	f, _ := os.Create(path)
	defer f.Close()
	pem.Encode(f, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
}

// writeKey writes PEM-encoded ECDSA private key to file
func writeKey(path string, key *ecdsa.PrivateKey) {
	f, _ := os.Create(path)
	defer f.Close()
	b, _ := x509.MarshalECPrivateKey(key)
	pem.Encode(f, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
}
