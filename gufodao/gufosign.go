package gufodao

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

// GufoSign sets the correct Sign value depending on security.mode
func Gufosign(t *pb.Request) *pb.Request {
	mode := strings.ToLower(viper.GetString("security.mode"))

	switch mode {

	// -----------------------------
	// STATIC SIGN MODE
	// -----------------------------
	case "sign":
		s := viper.GetString("server.sign")
		if s != "" {
			t.Sign = &s
		}

	// -----------------------------
	// HMAC MODE
	// -----------------------------
	case "hmac":
		secret := viper.GetString("security.hmac_secret")

		if t.Module == nil {
			empty := ""
			t.Module = &empty
		}

		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(*t.Module))
		sum := mac.Sum(nil)

		sign := hex.EncodeToString(sum)
		t.Sign = &sign

	// -----------------------------
	// MTLS MODE
	// -----------------------------
	case "mtls":
		// mTLS uses certificates only; Sign must be nil
		t.Sign = nil

	// -----------------------------
	// UNKNOWN MODE
	// -----------------------------
	default:
		// Fail-safe: remove sign
		t.Sign = nil
	}

	return t
}
