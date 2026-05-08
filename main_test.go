package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", res.StatusCode)
	}

	var body HealthResponse
	json.NewDecoder(res.Body).Decode(&body)

	if body.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", body.Status)
	}
}

func TestItemsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	w := httptest.NewRecorder()

	itemsHandler(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", res.StatusCode)
	}

	var body []Item
	json.NewDecoder(res.Body).Decode(&body)

	if len(body) != 3 {
		t.Errorf("expected 3 items, got %d", len(body))
	}
}