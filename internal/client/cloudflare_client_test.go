package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func cloudFlareClientTestSetup(t *testing.T) {
	t.Setenv("CLOUDFLAREDDNS_ZONE_ID", "test-zone-id")
	t.Setenv("CLOUDFLAREDDNS_API_TOKEN", "test-api-token")
}

func TestUpdateDNSRecordContent_Success(t *testing.T) {
	cloudFlareClientTestSetup(t)
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Validate incoming request details if needed
		if r.URL.Path != "/client/v4/zones/test-zone-id/dns_records/test-dns-record-id" {
			t.Errorf("Expected path '/client/v4/zones/test-zone-id/dns_records/test-dns-record-id', got '%s'", r.URL.Path)
		}
		if r.Method != "PATCH" {
			t.Errorf("Expected method 'PATCH', got '%s'", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-api-token" {
			t.Errorf("Expected Authorization header 'Bearer test-api-token', got '%s'", r.Header.Get("Authorization"))
		}
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Failed to read request body: %v", err)
		}
		expectedBody := `{"content":"1.2.3.4"}`
		if string(bodyBytes) != expectedBody {
			t.Errorf("Expected request body '%s', got '%s'", expectedBody, string(bodyBytes))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success": true}`))
	},
	))
	defer mockServer.Close()

	client := NewCloudFlareClient()
	// Use the mock server's client
	client.HttpClient = mockServer.Client()
	// Override the API URL to point to the mock server
	client.APIURL = mockServer.URL

	err := client.UpdateDNSRecordContent("test-dns-record-id", "1.2.3.4")
	if err != nil {
		t.Errorf("Expected UpdateDNSRecordContent to succeed, got: %v", err)
	}

}

func TestUpdateDNSRecordContent_EmptyDNSRecordID(t *testing.T) {
	cloudFlareClientTestSetup(t)

	err := NewCloudFlareClient().UpdateDNSRecordContent("", "1.2.3.4")
	if err == nil {
		t.Errorf("Expected UpdateDNSRecordContent to fail with empty DNS record ID, got: %v", err)
	}
}

func TestUpdateDNSRecordContent_EmptyNewContent(t *testing.T) {
	cloudFlareClientTestSetup(t)

	err := NewCloudFlareClient().UpdateDNSRecordContent("test-dns-record-id", "")
	if err == nil {
		t.Errorf("Expected UpdateDNSRecordContent to fail with empty new content, got: %v", err)
	}
}

func TestUpdateDNSRecordContent_FailedRequest(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)

	}))
	defer mockServer.Close()

	client := NewCloudFlareClient()
	// Use the mock server's client
	client.HttpClient = mockServer.Client()
	// Override the API URL to point to the mock server
	client.APIURL = mockServer.URL

	// Call the function under test
	err := client.UpdateDNSRecordContent("test-dns-record-id", "1.2.3.4")
	if err == nil {
		t.Fatalf("Expected UpdateDNSRecordContent to throw an error on failed request")
	}
}

func TestUpdateDNSRecordContent_UnexpectedStatusCode(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	client := NewCloudFlareClient()
	// Use the mock server's client
	client.HttpClient = mockServer.Client()
	// Override the API URL to point to the mock server
	client.APIURL = mockServer.URL

	// Call the function under test
	err := client.UpdateDNSRecordContent("test-dns-record-id", "1.2.3.4")
	if err == nil {
		t.Fatalf("Expected UpdateDNSRecordContent to throw an error on failed request")
	}
}
