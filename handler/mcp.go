package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	sf "github.com/gogufo/gufo-api-gateway/gufodao"
	pb "github.com/gogufo/gufo-api-gateway/proto/go"
	"github.com/gogufo/gufo-api-gateway/transport"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// MCP handles MCP HTTP requests and forwards them to the configured
// microservice through the existing Gufo gRPC transport.
//
// The MCP protocol itself is handled by the target microservice.
// Gufo only:
//   1. accepts the MCP HTTP request;
//   2. wraps the request into pb.Request;
//   3. marks it as MCP;
//   4. sends it through the existing gRPC transport;
//   5. extracts the MCP response;
//   6. returns the MCP JSON to the HTTP client.
func MCP(w http.ResponseWriter, r *http.Request) {
	mcpPath := viper.GetString("mcp.path")
	mcpModule := viper.GetString("mcp.module")

	if mcpPath == "" || mcpModule == "" {
		http.Error(w, "MCP is not configured", http.StatusNotFound)
		return
	}

	// MCP JSON-RPC requests are currently accepted via POST.
	if r.Method != http.MethodPost {
		http.Error(w, "MCP endpoint requires POST", http.StatusMethodNotAllowed)
		return
	}

	// Read the complete MCP JSON-RPC request.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read MCP request body", http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		http.Error(w, "Empty MCP request body", http.StatusBadRequest)
		return
	}

	// Validate JSON without unmarshalling and re-marshalling it.
	// The original request bytes are preserved below.
	if !json.Valid(body) {
		http.Error(w, "Invalid MCP JSON", http.StatusBadRequest)
		return
	}

	// Extract JSON-RPC request ID so that Gufo-generated errors
	// can preserve the request/response correlation.
	requestID := extractMCPRequestID(body)

	// Build the standard Gufo request.
	req := &pb.Request{
		Module: mcpModule,
		Path:   mcpPath,
		Method: pb.Method_METHOD_POST,
		Auth:   &pb.AuthContext{},
		Context: &pb.RequestContext{
			ApiVersion: "v1",
			Meta: map[string]string{
				"mcp": "true",
			},
		},
	}

	// Preserve the original MCP JSON bytes as protobuf Any.
	bytesValue := &wrapperspb.BytesValue{
		Value: body,
	}

	req.Body, err = anypb.New(bytesValue)
	if err != nil {
		http.Error(w, "Cannot encode MCP request", http.StatusInternalServerError)
		return
	}

	// Preserve authentication information exactly as for normal requests.
	req = fillAuthFromHeaders(req, r)

	// Preserve request metadata.
	if req.Context == nil {
		req.Context = &pb.RequestContext{}
	}

	if req.Context.Meta == nil {
		req.Context.Meta = make(map[string]string)
	}

	req.Context.Meta["mcp"] = "true"

	if rid := r.Header.Get("X-Request-ID"); rid != "" {
		req.Context.RequestId = rid
	}

	req.Context.Ip = r.RemoteAddr
	req.Context.UserAgent = r.UserAgent()

	// Use the existing Gufo gRPC transport.
	tr := transport.Get()

	resp, err := tr.Call(
		r.Context(),
		req.Module,
		sf.ProtoMethodToString(req.Method),
		req,
	)
	if err != nil {
		writeMCPError(
			w,
			http.StatusBadGateway,
			-32603,
			fmt.Sprintf("MCP transport error: %v", err),
			requestID,
		)
		return
	}

	if resp == nil {
		writeMCPError(
			w,
			http.StatusBadGateway,
			-32603,
			"Empty response from microservice",
			requestID,
		)
		return
	}

	// The MCP microservice must return the MCP JSON payload
	// in response.data["mcp"].
	mcpAny, ok := resp.Data["mcp"]
	if !ok || mcpAny == nil {
		writeMCPError(
			w,
			http.StatusBadGateway,
			-32603,
			"Microservice response does not contain MCP payload",
			requestID,
		)
		return
	}

	mcpJSON, err := extractMCPBody(mcpAny)
	if err != nil {
		writeMCPError(
			w,
			http.StatusBadGateway,
			-32603,
			fmt.Sprintf("Invalid MCP response from microservice: %v", err),
			requestID,
		)
		return
	}

	// Return the MCP JSON-RPC response to Hermes.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(mcpJSON)
}

// extractMCPRequestID extracts the JSON-RPC request ID from the original
// MCP request without changing or reserializing the request body.
func extractMCPRequestID(body []byte) interface{} {
	var request struct {
		ID interface{} `json:"id"`
	}

	if err := json.Unmarshal(body, &request); err != nil {
		return nil
	}

	return request.ID
}

// extractMCPBody extracts raw JSON bytes from the Any representation
// produced by the MCP microservice.
func extractMCPBody(v *anypb.Any) ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("nil protobuf Any")
	}

	// ConvertInterfaceToAny stores JSON bytes inside BytesValue.
	const bytesValueTypeURL = "type.googleapis.com/google.protobuf.BytesValue"

	if v.TypeUrl != bytesValueTypeURL {
		return nil, fmt.Errorf(
			"unexpected Any type: %s",
			v.TypeUrl,
		)
	}

	bytesValue := &wrapperspb.BytesValue{}

	if err := v.UnmarshalTo(bytesValue); err != nil {
		return nil, err
	}

	if len(bytesValue.Value) == 0 {
		return nil, fmt.Errorf("empty MCP response")
	}

	// Verify that the microservice returned valid JSON.
	if !json.Valid(bytesValue.Value) {
		return nil, fmt.Errorf("MCP response is not valid JSON")
	}

	return bytesValue.Value, nil
}

// writeMCPError returns a JSON-RPC error response.
func writeMCPError(
	w http.ResponseWriter,
	httpStatus int,
	code int,
	message string,
	id interface{},
) {
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	_, _ = w.Write(data)
}