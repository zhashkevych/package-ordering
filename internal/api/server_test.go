package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func performRequest(r http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestHealth(t *testing.T) {
	s := NewServer([]int{250, 500})
	w := performRequest(s.Router(), http.MethodGet, "/health", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestPacksGetAndSet(t *testing.T) {
	s := NewServer([]int{250, 1000, 500})

	// initial GET should return sorted ascending
	w := performRequest(s.Router(), http.MethodGet, "/packs", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var got map[string][]int
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(got["packSizes"]) != 3 || got["packSizes"][0] != 250 || got["packSizes"][1] != 500 || got["packSizes"][2] != 1000 {
		t.Fatalf("unexpected initial packSizes: %v", got["packSizes"])
	}

	// set new pack sizes
	setReq := SetPacksRequest{PackSizes: []int{23, 31, 53, 53}}
	w = performRequest(s.Router(), http.MethodPut, "/packs", setReq)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 on set packs, got %d", w.Code)
	}
	got = map[string][]int{}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode set response: %v", err)
	}
	// server stores descending unique order internally
	packs := got["packSizes"]
	if len(packs) != 3 || packs[0] != 53 || packs[1] != 31 || packs[2] != 23 {
		t.Fatalf("unexpected stored packSizes: %v", packs)
	}
}

func TestCalculateEndpoint(t *testing.T) {
	s := NewServer([]int{23, 31, 53})

	calcReq := CalculateRequest{Amount: 500000}
	w := performRequest(s.Router(), http.MethodPost, "/calculate", calcReq)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp CalculateResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode calc response: %v", err)
	}
	if resp.Amount != 500000 {
		t.Fatalf("unexpected amount: %d", resp.Amount)
	}
	if resp.Allocation[23] != 2 || resp.Allocation[31] != 7 || resp.Allocation[53] != 9429 {
		t.Fatalf("unexpected allocation: %+v", resp.Allocation)
	}
	if resp.TotalItems < resp.Amount {
		t.Fatalf("total items must be >= amount, got %d < %d", resp.TotalItems, resp.Amount)
	}
}

func TestCalculateBadRequests(t *testing.T) {
	s := NewServer([]int{250, 500})

	// invalid json
	req := httptest.NewRequest(http.MethodPost, "/calculate", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Router().ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	// amount <= 0
	w = performRequest(s.Router(), http.MethodPost, "/calculate", CalculateRequest{Amount: 0})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for non-positive amount, got %d", w.Code)
	}
}
