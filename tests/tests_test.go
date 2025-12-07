// Copyright 2020-2025 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
