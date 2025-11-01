package gufodao

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	viper "github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// --- Sign / HMAC verification ---

// ComputeHMAC generates an HMAC-based signature using secret, module, and timestamp.
func ComputeHMAC(secret, module string, ts int64) string {
	data := fmt.Sprintf("%s:%d", module, ts)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

// VerifyHMAC validates the HMAC signature.
func VerifyHMAC(secret, module, sign string, maxAge time.Duration) bool {
	parts := strings.Split(sign, ":")
	if len(parts) != 2 {
		return false
	}
	sig := parts[0]
	ts, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return false
	}
	if time.Since(time.Unix(ts, 0)) > maxAge {
		return false
	}
	expected := ComputeHMAC(secret, module, ts)
	return hmac.Equal([]byte(sig), []byte(expected))
}

// --- TLS credentials (for mTLS mode) ---

// LoadMTLSCredentials loads client-side TLS credentials for gRPC connections.
func LoadMTLSCredentials() (credentials.TransportCredentials, error) {
	caCertPath := viper.GetString("security.ca_cert")
	clientCertPath := viper.GetString("security.cert")
	clientKeyPath := viper.GetString("security.key")

	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}

	cert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load client key pair: %w", err)
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("failed to append CA certs")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	return credentials.NewTLS(tlsConfig), nil
}

// --- Utility: GetGRPCCredentials chooses credentials based on mode ---

// GetGRPCCredentials returns proper transport credentials depending on security mode.
func GetGRPCCredentials() (grpc.DialOption, error) {
	mode := strings.ToLower(viper.GetString("security.mode"))

	switch mode {
	case "mtls":
		creds, err := LoadMTLSCredentials()
		if err != nil {
			return nil, err
		}
		return grpc.WithTransportCredentials(creds), nil

	default:
		// for sign and hmac modes, use insecure (plain gRPC)
		return grpc.WithInsecure(), nil
	}
}
