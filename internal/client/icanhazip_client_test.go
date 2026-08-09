package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPublicIP_Success(t *testing.T) {

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate incoming request details if needed
		if r.URL.Path != "/" {
			t.Errorf("Expected path '/', got '%s'", r.URL.Path)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("1.2.3.4"))
	}))
	defer mockServer.Close()

	client := NewICanHazIPClient()
	// Use the mock server's client
	client.HttpClient = mockServer.Client()
	// Override the API URL to point to the mock server
	client.APIURL = mockServer.URL

	// Call the function under test
	ip, err := client.GetPublicIP()
	if err != nil {
		t.Fatalf("Expected GetPublicIP to succeed, got: %v", err)
	}
	if ip != "1.2.3.4" {
		t.Errorf("Expected returned IP to be '1.2.3.4', got: %s", ip)
	}
}
