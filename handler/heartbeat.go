package handler

import (
	"encoding/json"
	"net/http"
	"time"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/spf13/viper"
)

// Universal heartbeat entry through Gateway.
// Microservices POST here -> Gufo routes to masterservice.
func HeartbeatHandler(w http.ResponseWriter, r *http.Request, t *pb.Request) {

	msEnabled := viper.GetBool("server.masterservice")

	// ------------------------------------------------------------
	// MODE 2: Standalone → mock MasterService
	// ------------------------------------------------------------
	if !msEnabled {
		mock := map[string]interface{}{
			"leader": true,
			"cron":   true,
			"ttl":    0,
			"epoch":  0,
			"ts":     time.Now().Unix(),
		}
		moduleAnswerv3(w, r, mock, t)
		return
	}

	// ------------------------------------------------------------
	// MODE 1: Cluster → proxy to MasterService
	// ------------------------------------------------------------
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		errorAnswer(w, r, t, 400, "0000500", "Invalid JSON body")
		return
	}
	defer r.Body.Close()

	payload["ts"] = time.Now().Unix()

	req := &pb.Request{
		Module: sf.StringPtr("masterservice"),
		IR: &pb.InternalRequest{
			Param:  sf.StringPtr("heartbeat"),
			Method: sf.StringPtr("POST"),
		},
		Args: sf.ToMapStringAny(payload),
	}

	ans := sf.GRPCConnect(
		sf.ConfigString("microservices.masterservice.host"),
		sf.ConfigString("microservices.masterservice.port"),
		req,
	)

	if ans == nil || ans["httpcode"] != nil {
		errorAnswer(w, r, t, 500, "0000501", "MasterService heartbeat error")
		return
	}

	moduleAnswerv3(w, r, ans, t)
}
