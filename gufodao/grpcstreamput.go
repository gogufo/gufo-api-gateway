package gufodao

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	viper "github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/anypb"
)

// GRPCStreamPut streams either a single raw-body file (with X-Filename) or multiple files via multipart/form-data.
// It sends FileChunk messages wrapped into Request.Args["chunk"] (google.protobuf.Any).
func GRPCStreamPut(host, port string, r *http.Request, t *pb.Request) map[string]interface{} {

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
	// IMPORTANT: r.Body закрываем здесь; gRPC conn — из пула, НЕ закрываем.

	timeout := viper.GetDuration(fmt.Sprintf("microservices.%s.stream_timeout", module))
	if timeout == 0 {
		timeout = 2 * time.Minute
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := pb.NewReverseClient(conn)
	stream, err := client.Stream(ctx)
	if err != nil {
		return map[string]interface{}{"httpcode": 500, "message": fmt.Sprintf("open stream failed: %v", err)}
	}
	defer stream.CloseSend()
	defer r.Body.Close()

	ct := r.Header.Get("Content-Type")
	if strings.HasPrefix(strings.ToLower(ct), "multipart/") {
		if err := streamMultipartFiles(stream, r, t); err != nil {
			return map[string]interface{}{"httpcode": 500, "message": err.Error()}
		}
	} else {
		if err := streamSingleBody(stream, r, t); err != nil {
			return map[string]interface{}{"httpcode": 500, "message": err.Error()}
		}
	}

	// Receive server responses (optional aggregation)
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return map[string]interface{}{"httpcode": 500, "message": fmt.Sprintf("recv error: %v", err)}
		}
		// Можно логировать resp.Data, если нужно
		_ = resp
	}

	return map[string]interface{}{"httpcode": 200, "message": "PUT streamed successfully"}
}

func streamSingleBody(stream pb.Reverse_StreamClient, r *http.Request, t *pb.Request) error {
	filename := r.Header.Get("X-Filename")
	if filename == "" {
		filename = "upload.bin"
	}
	buf := make([]byte, 64*1024)

	// Отправляем «start-of-file» маркер
	if err := sendFileMarker(stream, t, filename, "start"); err != nil {
		return err
	}

	for {
		n, err := r.Body.Read(buf)
		if n > 0 {
			if err := sendChunk(stream, t, filename, buf[:n]); err != nil {
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

	// «end-of-file» маркер
	return sendFileMarker(stream, t, filename, "end")
}

func streamMultipartFiles(stream pb.Reverse_StreamClient, r *http.Request, t *pb.Request) error {
	mr, err := r.MultipartReader()
	if err != nil {
		return fmt.Errorf("multipart reader: %w", err)
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("multipart next part: %w", err)
		}
		if part.FileName() == "" {
			// это не файл (form-field) — пропускаем
			continue
		}

		filename := part.FileName()
		buf := make([]byte, 64*1024)

		// start-of-file
		if err := sendFileMarker(stream, t, filename, "start"); err != nil {
			return err
		}

		for {
			n, err := part.Read(buf)
			if n > 0 {
				if err := sendChunk(stream, t, filename, buf[:n]); err != nil {
					return err
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return fmt.Errorf("read part %s: %w", filename, err)
			}
		}

		// end-of-file
		if err := sendFileMarker(stream, t, filename, "end"); err != nil {
			return err
		}

		part.Close()
	}

	return nil
}

func sendChunk(stream pb.Reverse_StreamClient, t *pb.Request, filename string, data []byte) error {
	chunk := &pb.FileChunk{Name: filename, Data: data}
	anyChunk, err := anypb.New(chunk)
	if err != nil {
		return err
	}
	req := &pb.Request{
		Module: t.Module,
		IR:     t.IR,
		Args:   map[string]*anypb.Any{"chunk": anyChunk},
	}
	return stream.Send(req)
}

// sendFileMarker — служебный маркер начала/конца файла (FileChunk с пустыми данными и именем + флагом)
func sendFileMarker(stream pb.Reverse_StreamClient, t *pb.Request, filename, phase string) error {
	meta := map[string]string{"filename": filename, "phase": phase} // {"start","end"}
	anyMeta, err := anypb.New(&pb.StringMap{Entries: meta})
	if err != nil {
		return err
	}
	req := &pb.Request{
		Module: t.Module,
		IR:     t.IR,
		Args:   map[string]*anypb.Any{"meta": anyMeta},
	}
	return stream.Send(req)
}
