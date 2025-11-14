package tests

import (
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	resp, err := http.Get("http://localhost:8090/api/v3/health")
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("health check failed: %v (%d)", err, resp.StatusCode)
	}
}

func TestMetrics(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost:9100/api/v3/metrics", nil)
	req.Header.Set("X-Metrics-Token", "gufo-metrics")
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("metrics endpoint failed: %v (%d)", err, resp.StatusCode)
	}
}
