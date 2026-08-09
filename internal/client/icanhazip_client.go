package client

import (
	"fmt"
	"io"
	"net/http"
)

type ICanHazIPClient struct {
	HttpClient *http.Client
	APIURL     string
}

func NewICanHazIPClient() *ICanHazIPClient {
	return &ICanHazIPClient{
		HttpClient: http.DefaultClient,
		APIURL:     "https://ipv4.icanhazip.com",
	}
}

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
