// grpcstream.go
package gufodao

import (
	"context"
	"fmt"
	"io"
	"time"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
)

// GRPCStream handles bidirectional streaming calls
func GRPCStream(host, port string, t *pb.Request) map[string]interface{} {
	result := make(map[string]interface{})

	conn, err := GetGRPCConn(
		host, port,
		"certs/ca.crt", "certs/client.crt", "certs/client.key",
		false, // useMTLS = false for now
	)
	if err != nil {
		SetErrorLog(fmt.Sprintf("stream: connect failed %v", err))
		result["error"] = err.Error()
		return result
	}
	client := pb.NewReverseClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	stream, err := client.Stream(ctx)
	if err != nil {
		SetErrorLog("stream: open failed " + err.Error())
		result["error"] = err.Error()
		return result
	}

	if err := stream.Send(t); err != nil {
		SetErrorLog("stream: send failed " + err.Error())
		return result
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			SetErrorLog("stream: recv failed " + err.Error())
			break
		}
		fmt.Println("STREAM DATA:", resp.Data)
	}
	stream.CloseSend()

	result["status"] = "ok"
	return result
}
