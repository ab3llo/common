package common

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	recorder := httptest.NewRecorder()
	data := map[string]string{"foo": "bar"}
	WriteJSON(recorder, http.StatusOK, data)

	if recorder.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if ct := recorder.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
	var resp map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}
	if resp["foo"] != "bar" {
		t.Errorf("expected foo=bar, got %v", resp)
	}
}

func TestWriteJSON_EncodeError(t *testing.T) {
	recorder := httptest.NewRecorder()
	ch := make(chan int) // channels can't be JSON encoded
	WriteJSON(recorder, http.StatusInternalServerError, ch)
	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, recorder.Code)
	}
}

func TestReadJSON(t *testing.T) {
	data := map[string]string{"foo": "bar"}
	body, _ := json.Marshal(data)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	var out map[string]string
	err := ReadJSON(req, &out)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if out["foo"] != "bar" {
		t.Errorf("expected foo=bar, got %v", out)
	}
}

func TestReadJSON_Invalid(t *testing.T) {
	req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte("notjson")))
	var out map[string]string
	err := ReadJSON(req, &out)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestWriteError(t *testing.T) {
	recorder := httptest.NewRecorder()
	WriteError(recorder, http.StatusBadRequest, "something went wrong")
	if recorder.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
	var resp map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &resp); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}
	if resp["error"] != "something went wrong" {
		t.Errorf("expected error message, got %v", resp)
	}
}

func BenchmarkWriteJSON(b *testing.B) {
	data := map[string]string{"foo": "bar", "baz": "qux"}
	for i := 0; i < b.N; i++ {
		recorder := httptest.NewRecorder()
		WriteJSON(recorder, http.StatusOK, data)
	}
}

func BenchmarkReadJSON(b *testing.B) {
	data := map[string]string{"foo": "bar"}
	body, _ := json.Marshal(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		var out map[string]string
		ReadJSON(req, &out)
	}
}
