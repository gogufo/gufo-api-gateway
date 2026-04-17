package handler

import (
	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

// Handle is the replacement for InternalRequest
func Handle(t *pb.Request) *pb.Response {

	// --- heartbeat shortcut ---
	if t.Module == "heartbeat" {
		ans, err := heartbeatCore(t, nil)
		if err != nil {
			return sf.ErrorReturn(t, 500, "0000501", "MasterService heartbeat error")
		}
		return sf.Interfacetoresponse(t, ans)
	}

	// --- resolve microservice ---
	host, port, _ := GetHostAndPort(t)

	// --- gRPC proxy ---
	ans := sf.GRPCConnect(host, port, t)

	// --- convert response ---
	return sf.Interfacetoresponse(t, ans)
}
