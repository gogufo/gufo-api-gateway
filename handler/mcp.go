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

func MCP(w http.ResponseWriter, r *http.Request) {
	mcpPath := viper.GetString("mcp.path")
	mcpModule := viper.GetString("mcp.module")

	if mcpPath == "" || mcpModule == "" {
		http.Error(w, "MCP is not configured", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "MCP endpoint requires POST", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read MCP request body", http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		http.Error(w, "Empty MCP request body", http.StatusBadRequest)
		return
	}

	if !json.Valid(body) {
		http.Error(w, "Invalid MCP JSON", http.StatusBadRequest)
		return
	}

	requestID := extractMCPRequestID(body)

	// Build the standard Gufo request first.
	// This preserves the same Auth, Sign, IP, User-Agent and
	// other request initialization used by normal REST requests.
	req := RequestInit(r)

	// MCP-specific routing.
	req.Module = mcpModule
	req.Path = mcpPath
	req.Method = pb.Method_METHOD_POST

	// Mark the request as MCP for the target microservice.
	if req.Context == nil {
		req.Context = &pb.RequestContext{}
	}

	req.Context.ApiVersion = "v1"

	if req.Context.Meta == nil {
		req.Context.Meta = make(map[string]string)
	}

	req.Context.Meta["mcp"] = "true"

	// Preserve the HTTP request ID explicitly.
	if rid := r.Header.Get("X-Request-ID"); rid != "" {
		req.Context.RequestId = rid
	}

	// Store raw MCP JSON in protobuf Any.
	bytesValue := &wrapperspb.BytesValue{
		Value: body,
	}

	req.Body, err = anypb.New(bytesValue)
	if err != nil {
		http.Error(w, "Cannot encode MCP request", http.StatusInternalServerError)
		return
	}

	// Preserve normal Gufo authentication headers.
	req = fillAuthFromHeaders(req, r)

	// Send the request through the existing Gufo transport.
	tr := transport.Get()

	fmt.Println(">>> GUFO MCP: calling microservice")
	fmt.Printf(">>> GUFO MCP: module=%s method=%s path=%s sign_set=%t\n",
		req.Module,
		sf.ProtoMethodToString(req.Method),
		req.Path,
		req.Auth != nil && req.Auth.Sign != "",
	)

	resp, err := tr.Call(
		r.Context(),
		req.Module,
		sf.ProtoMethodToString(req.Method),
		req,
	)

	fmt.Printf(">>> GUFO MCP: response=%+v err=%v\n", resp, err)

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(mcpJSON)
}

func extractMCPRequestID(body []byte) json.RawMessage {
	var request struct {
		ID json.RawMessage `json:"id"`
	}

	if err := json.Unmarshal(body, &request); err != nil {
		return nil
	}

	if len(request.ID) == 0 {
		return nil
	}

	return request.ID
}

func extractMCPBody(v *anypb.Any) ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("nil protobuf Any")
	}

	const bytesValueTypeURL = "type.googleapis.com/google.protobuf.BytesValue"

	if v.TypeUrl != bytesValueTypeURL {
		return nil, fmt.Errorf("unexpected Any type: %s", v.TypeUrl)
	}

	bytesValue := &wrapperspb.BytesValue{}

	if err := v.UnmarshalTo(bytesValue); err != nil {
		return nil, err
	}

	if len(bytesValue.Value) == 0 {
		return nil, fmt.Errorf("empty MCP response")
	}

	if !json.Valid(bytesValue.Value) {
		return nil, fmt.Errorf("MCP response is not valid JSON")
	}

	return bytesValue.Value, nil
}

func writeMCPError(
	w http.ResponseWriter,
	httpStatus int,
	code int,
	message string,
	id json.RawMessage,
) {
	if len(id) == 0 {
		id = json.RawMessage("null")
	}

	response := struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      json.RawMessage `json:"id"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}{
		JSONRPC: "2.0",
		ID:      id,
	}

	response.Error.Code = code
	response.Error.Message = message

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_, _ = w.Write(data)
}