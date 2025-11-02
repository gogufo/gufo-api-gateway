package gufodao

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	viper "github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/anypb"
)

// GRPCStream handles bidirectional streaming (e.g. multi-file upload)
func GRPCStream(host, port string, t *pb.Request) map[string]interface{} {
	answer := make(map[string]interface{})
	addr := fmt.Sprintf("%s:%s", host, port)

	conn, err := GetGRPCConn(
		host,
		port,
		viper.GetString("security.ca_path"),
		viper.GetString("security.cert_path"),
		viper.GetString("security.key_path"),
		viper.GetBool("security.mtls"),
	)
	if err != nil {
		answer["httpcode"] = 400
		answer["message"] = fmt.Sprintf("Failed to connect to %s: %v", addr, err)
		return answer
	}

	client := pb.NewReverseClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	stream, err := client.Stream(ctx)
	if err != nil {
		answer["httpcode"] = 500
		answer["message"] = fmt.Sprintf("Failed to open stream: %v", err)
		return answer
	}

	// List of files to upload
	files := t.Files // добавь []string в pb.Request (список путей)
	for _, path := range files {
		if err := sendFile(stream, path); err != nil {
			answer["httpcode"] = 500
			answer["message"] = fmt.Sprintf("upload %s failed: %v", path, err)
			return answer
		}
	}

	// Signal EOF to server
	stream.CloseSend()

	// Collect responses from server
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			answer["httpcode"] = 500
			answer["message"] = fmt.Sprintf("recv error: %v", err)
			return answer
		}

		fmt.Println("📦 Server responded:", resp.Data)
	}

	answer["httpcode"] = 200
	answer["message"] = "All files uploaded successfully"
	return answer
}

// sendFile sends one file as chunks via stream
func sendFile(stream pb.Reverse_StreamClient, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, 64*1024)
	filename := filepath.Base(path)

	for {
		n, err := file.Read(buf)
		if n > 0 {
			chunk, _ := anypb.New(&pb.FileChunk{
				Name: filename,
				Data: buf[:n],
			})
			req := &pb.Request{
				Args: map[string]*anypb.Any{"chunk": chunk},
			}
			if err := stream.Send(req); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}
