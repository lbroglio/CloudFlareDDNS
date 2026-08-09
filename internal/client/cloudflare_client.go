package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/lbroglio/CloudFlareDDNS/internal/config"
)

type CloudFlareClient struct {
	HttpClient *http.Client
	APIURL     string
}

// NewCloudFlareClient creates a new instance of CloudFlareClient with default settings.
// It uses the default HTTP client and the URL https://api.cloudflare.com.
func NewCloudFlareClient() *CloudFlareClient {
	return &CloudFlareClient{
		HttpClient: http.DefaultClient,
		APIURL:     "https://api.cloudflare.com",
	}
}

func (c *CloudFlareClient) buildRequestForURL(path string, requestType string) (*http.Request, error) {
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}
	if requestType != "GET" && requestType != "POST" && requestType != "PUT" && requestType != "PATCH" && requestType != "DELETE" {
		return nil, fmt.Errorf("requestType must be one of GET, POST, PUT, PATCH, DELETE")
	}

	req, err := http.NewRequest(requestType, c.APIURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %v", err)
	}

	apiToken := config.GetCloudflareAPIToken()
	req.Header.Set("Authorization", "Bearer "+apiToken)

	return req, nil
}

func (c *CloudFlareClient) UpdateDNSRecordContent(dnsRecordID string, newContent string) error {
	if dnsRecordID == "" {
		return fmt.Errorf("dnsRecordID cannot be empty")
	}
	if newContent == "" {
		return fmt.Errorf("newContent cannot be empty")
	}

	req, err := c.buildRequestForURL("/client/v4/zones/"+config.GetCloudflareZoneID()+"/dns_records/"+dnsRecordID, "PATCH")
	if err != nil {
		return fmt.Errorf("failed to build request: %v", err)
	}

	body := "{\"content\":\"" + newContent + "\"}"
	bodyBuffer := io.NopCloser(bytes.NewBuffer([]byte(body)))
	// Close the buffer after use to avoid memory leaks
	defer bodyBuffer.Close()
	req.Body = bodyBuffer

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}

	// Close the response body after use to avoid memory leaks
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code for %s request: %d", c.APIURL, resp.StatusCode)
	}

	return nil
}
