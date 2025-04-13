package apiclient

import (
	"CountrySearch/model"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func mockServer(t *testing.T, handlerFunc http.HandlerFunc) (*httptest.Server, func()) {
	server := httptest.NewServer(handlerFunc)

	// Replace the global URL with mock server URL
	originalURL := URL
	URL = server.URL + "/"
	cleanup := func() {
		server.Close()
		URL = originalURL
	}
	return server, cleanup
}

func TestFetchCountryDataFromAPI_Success(t *testing.T) {
	mockJSON := []map[string]interface{}{
		{
			"name": map[string]interface{}{
				"common": "India",
			},
			"capital":    []interface{}{"New Delhi"},
			"currencies": map[string]interface{}{"INR": map[string]interface{}{"symbol": "₹"}},
			"population": 1380004385,
		},
	}
	mockResp, _ := json.Marshal(mockJSON)

	_, cleanup := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(mockResp)
	})
	defer cleanup()

	client := New()
	ctx := context.Background()
	resp, err := client.FetchCountryDataFromAPI(ctx, "India")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	expected := model.Response{
		Name:       "India",
		Capital:    "New Delhi",
		Currency:   "₹",
		Population: 1380004385,
	}
	if resp != expected {
		t.Errorf("Unexpected response: %+v", resp)
	}
}

func TestFetchCountryDataFromAPI_RequestError(t *testing.T) {
	// Invalid context to trigger request error
	client := New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.FetchCountryDataFromAPI(ctx, "India")
	if err == nil || !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("Expected context error, got: %v", err)
	}
}

func TestFetchCountryDataFromAPI_StatusCodeError(t *testing.T) {
	_, cleanup := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Forbidden", http.StatusForbidden)
	})
	defer cleanup()

	client := New()
	ctx := context.Background()
	_, err := client.FetchCountryDataFromAPI(ctx, "India")
	if err == nil || !strings.Contains(err.Error(), "Oops!! status code") {
		t.Errorf("Expected status code error, got: %v", err)
	}
}

func TestFetchCountryDataFromAPI_BadJSON(t *testing.T) {
	_, cleanup := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not-json"))
	})
	defer cleanup()

	client := New()
	ctx := context.Background()
	_, err := client.FetchCountryDataFromAPI(ctx, "India")
	if err == nil || !strings.Contains(err.Error(), "error decoding") {
		t.Errorf("Expected JSON decode error, got: %v", err)
	}
}

func TestFetchCountryDataFromAPI_EmptyResponse(t *testing.T) {
	emptyResp, _ := json.Marshal([]interface{}{})

	_, cleanup := mockServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(emptyResp)
	})
	defer cleanup()

	client := New()
	ctx := context.Background()
	_, err := client.FetchCountryDataFromAPI(ctx, "India")
	if err == nil || !strings.Contains(err.Error(), "error decoding") {
		t.Errorf("Expected error for empty response, got: %v", err)
	}
}
