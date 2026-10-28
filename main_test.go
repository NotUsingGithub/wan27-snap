package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestURLsAndValidation(t *testing.T) {
	if WebsiteURL != "https://wan2-7.io/" {
		t.Fatalf("unexpected WebsiteURL: %s", WebsiteURL)
	}
	if err := ValidateID("wan2.7-t2v", "model"); err != nil {
		t.Fatalf("expected valid ID: %v", err)
	}
	if err := ValidateID("../etc/passwd", "model"); err == nil {
		t.Fatal("expected error for path traversal")
	}
}

func TestClientRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":      "req_123",
			"status":  "succeeded",
			"credits": 15,
		})
	}))
	defer srv.Close()

	c := NewClient("test-token", srv.URL)
	res, err := c.Request("POST", "/api/v1/generations", nil, map[string]any{"prompt": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["id"] != "req_123" {
		t.Fatalf("unexpected id: %v", res["id"])
	}
}
