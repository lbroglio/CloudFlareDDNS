// Package repsonsible for fetching configuration values from the environment variables.
//
// This packkage contains isolated functions which other packages can use to get
// configurations without needing to be coupled to the environment variable names directly.
//
// The reason to use environment variables for config is to make it easier to run this
// application in a containerized environment.
package config

import (
	"os"
	"strings"
)

// Get the root folder where the application will store its persistent data.
func GetPersistenceRootPath() string {
	return os.Getenv("CLOUDFLAREDDNS_PERSISTENCE_ROOT")
}

// Get the Cloudflare API token from the environment variable.
func GetCloudflareAPIToken() string {
	return os.Getenv("CLOUDFLAREDDNS_API_TOKEN")
}

// Get the Cloudflare Zone ID from the environment variable.
func GetCloudflareZoneID() string {
	return os.Getenv("CLOUDFLAREDDNS_ZONE_ID")
}

// GetTargetDNSRecordIDs retrieves the target DNS record IDs from the environment variable.
// It returns a slice of strings containing the DNS record IDs.
func GetTargetDNSRecordIDs() []string {
	// Get the comma-separated list of DNS record IDs from the environment variable.
	dnsRecordIDs := os.Getenv("CLOUDFLAREDDNS_DNS_RECORD_IDS")
	if dnsRecordIDs == "" {
		return []string{}
	}

	// Split the string into a slice of strings.
	targetRecordIDs := strings.Split(dnsRecordIDs, ",")
	for i := range targetRecordIDs {
		targetRecordIDs[i] = strings.TrimSpace(targetRecordIDs[i])
	}
	return targetRecordIDs
}
