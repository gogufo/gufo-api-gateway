package gufodao

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	viper "github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/anypb"
)

func GRPCStreamPut(host, port string, r *http.Request, t *pb.Request) map[string]interface{} {
	answer := make(map[string]interface{})
	module := safeModuleName(t)

	conn, err := GetGRPCConn(
		host, port,
		viper.GetString("security.ca_path"),
		viper.GetString("security.cert_path"),
		viper.GetString("security.key_path"),
		viper.GetBool("security.mtls"),
	)
	if err != nil {
		return map[string]interface{}{"httpcode": 400, "message": err.Error()}
	}
	defer r.Body.Close()

	timeout := viper.GetDuration(fmt.Sprintf("microservices.%s.stream_timeout", module))
	if timeout == 0 {
		timeout = 2 * time.Minute
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := pb.NewReverseClient(conn)
	stream, err := client.Stream(ctx)
	if err != nil {
		return map[string]interface{}{"httpcode": 500, "message": err.Error()}
	}
	defer stream.CloseSend()

	filename := r.Header.Get("X-Filename")
	if filename == "" {
		filename = "upload.bin"
	}

	buf := make([]byte, 64*1024)
	for {
		n, err := r.Body.Read(buf)
		if n > 0 {
			chunk := &pb.FileChunk{Name: filename, Data: buf[:n]}
			anyChunk, _ := anypb.New(chunk)
			req := &pb.Request{
				Module: t.Module,
				Args:   map[string]*anypb.Any{"chunk": anyChunk},
				IR:     t.IR,
			}
			if err := stream.Send(req); err != nil {
				return map[string]interface{}{"httpcode": 500, "message": err.Error()}
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return map[string]interface{}{"httpcode": 500, "message": err.Error()}
		}
	}

	// финальный пустой chunk как EOF-сигнал
	final := &pb.FileChunk{Name: filename, Data: []byte{}}
	finalAny, _ := anypb.New(final)
	stream.Send(&pb.Request{
		Module: t.Module,
		Args:   map[string]*anypb.Any{"chunk": finalAny},
		IR:     t.IR,
	})

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return map[string]interface{}{"httpcode": 500, "message": err.Error()}
		}
		fmt.Println("Server:", resp.Data)
	}

	answer["httpcode"] = 200
	answer["message"] = "PUT streamed successfully"
	return answer
}
