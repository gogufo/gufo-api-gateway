package handler

import (
	"encoding/json"
	"net/http"
	"time"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

// Universal heartbeat entry through Gateway.
// Microservices POST here -> Gufo routes to masterservice.
func HeartbeatHandler(w http.ResponseWriter, r *http.Request, t *pb.Request) {
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		errorAnswer(w, r, t, 400, "0000500", "Invalid JSON body")
		return
	}
	defer r.Body.Close()

	payload["ts"] = time.Now().Unix() // enrich with current timestamp

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

	moduleAnswerv3(w, r, ans, t)
}
