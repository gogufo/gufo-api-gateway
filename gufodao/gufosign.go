package gufodao

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

// GufoSign sets signature inside AuthContext according to security.mode
func Gufosign(t *pb.Request) *pb.Request {
	mode := strings.ToLower(viper.GetString("security.mode"))

	// Ensure AuthContext exists
	if t.Auth == nil {
		t.Auth = &pb.AuthContext{}
	}

	switch mode {

	// -----------------------------
	// STATIC SIGN MODE
	// -----------------------------
	case "sign":
		s := viper.GetString("server.sign")
		t.Auth.Sign = s

	// -----------------------------
	// HMAC MODE
	// -----------------------------
	case "hmac":
		secret := viper.GetString("security.hmac_secret")

		module := t.Module
		if module == "" {
			module = ""
		}

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(module))

		t.Auth.Sign = hex.EncodeToString(mac.Sum(nil))

	// -----------------------------
	// MTLS MODE
	// -----------------------------
	case "mtls":
		// mTLS → no signature
		t.Auth.Sign = ""

	// -----------------------------
	// UNKNOWN MODE
	// -----------------------------
	default:
		t.Auth.Sign = ""
	}

	return t
}
