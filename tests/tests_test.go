package tests

import (
	"net/http"
	"testing"
)

func TestHealth(t *testing.T) {
	resp, err := http.Get("http://localhost:8090/api/v3/health")
	if err != nil {
		t.Fatalf("health request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("health returned non-200 status: %d", resp.StatusCode)
	}
}

func TestMetrics(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:9100/api/v3/metrics", nil)
	if err != nil {
		t.Fatalf("cannot create metrics request: %v", err)
	}

	req.Header.Set("X-Metrics-Token", "gufo-metrics")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("metrics request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("metrics returned non-200 status: %d", resp.StatusCode)
	}
}
