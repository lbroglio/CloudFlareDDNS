package client

import (
	"fmt"
	"io"
	"net/http"
)

// ICanHazIPClient is a client for fetching the public IP address of the machine using the ICanHazIP API.
type ICanHazIPClient struct {
	HttpClient *http.Client
	APIURL     string
}

// NewICanHazIPClient creates a new instance of ICanHazIPClient with default settings.
// It uses the default HTTP client and the URL https://ipv4.icanhazip.com.
func NewICanHazIPClient() *ICanHazIPClient {
	return &ICanHazIPClient{
		HttpClient: http.DefaultClient,
		APIURL:     "https://ipv4.icanhazip.com",
	}
}

// GetPublicIP fetches the public IP address of the machine by making a GET request to the ICanHazIP API.
// It returns the public IP address as a string and an error if any occurred during the request or response processing.
func (c *ICanHazIPClient) GetPublicIP() (string, error) {
	resp, err := c.HttpClient.Get(c.APIURL)
	if err != nil {
		return "", fmt.Errorf("failed to get public IP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code for %s request: %d", c.APIURL, resp.StatusCode)
	}

	repsonseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(repsonseBytes), nil
}
