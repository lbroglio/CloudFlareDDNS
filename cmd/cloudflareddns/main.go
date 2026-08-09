package main

import (
	"fmt"
	"strings"

	"github.com/lbroglio/CloudFlareDDNS/internal/client"
	"github.com/lbroglio/CloudFlareDDNS/internal/config"
	"github.com/lbroglio/CloudFlareDDNS/internal/persistence"
)

func getIdsOfRecordsThatNeedUpdating(targetDNSRecordIDs []string, lastKnownIPs map[string]string, currentPublicIP string, cloudflareClient *client.CloudFlareClient) []string {
	// Slice to hold the IDs of DNS records that need updating
	recordsNeedUpdating := make([]string, 0, 10)
	for _, recordID := range targetDNSRecordIDs {
		lastKnownIP, exists := lastKnownIPs[recordID]
		if !exists || lastKnownIP == "" {
			// If there's no last known IP for this record check the current IP from Cloudflare
			dnsRecordDetails, err := cloudflareClient.GetDNSRecordDetails(recordID)
			if err != nil {
				panic(err)
			}
			lastKnownIP = dnsRecordDetails.Content
			// Update the last known IP in the map for future reference
			lastKnownIPs[recordID] = lastKnownIP
		}
		lastKnownIP = strings.TrimSpace(lastKnownIP)
		currentPublicIP = strings.TrimSpace(currentPublicIP)

		if lastKnownIP != currentPublicIP {
			recordsNeedUpdating = append(recordsNeedUpdating, recordID)
			persistence.WriteLogLine(fmt.Sprintf("DNS record %s needs updating. Last known IP: %s, Current public IP: %s", recordID, lastKnownIP, currentPublicIP))
		} else {
			persistence.WriteLogLine(fmt.Sprintf("DNS record %s does not need updating. Current IP: %s", recordID, currentPublicIP))
		}

	}

	return recordsNeedUpdating
}

func main() {
	targetDNSRecordIDs := config.GetTargetDNSRecordIDs()

	if len(targetDNSRecordIDs) == 0 {
		panic("No target DNS record IDs found in configuration. Please provide at least one DNS record ID to update.")
	}

	// Load the cached public
	recordCache, err := persistence.GetState()
	if err != nil {
		panic(err)
	}

	// Get the current public IP address
	hazipClient := client.NewICanHazIPClient()
	publicIP, err := hazipClient.GetPublicIP()
	if err != nil {
		panic(err)
	}
	publicIP = strings.TrimSpace(publicIP)

	cloudflareClient := client.NewCloudFlareClient()
	recordsToUpdate := getIdsOfRecordsThatNeedUpdating(targetDNSRecordIDs, recordCache.LastKnownIPsForDNS, publicIP, cloudflareClient)

	// Update the DNS records that need updating
	for _, recordID := range recordsToUpdate {
		persistence.WriteLogLine(fmt.Sprintf("Updating DNS record %s to new IP: %s", recordID, publicIP))
		err := cloudflareClient.UpdateDNSRecordContent(recordID, publicIP)
		if err != nil {
			panic(err)
		}
		recordCache.LastKnownIPsForDNS[recordID] = publicIP
	}

	// Save the updated cache
	persistence.UpdateState(recordCache)
	defer persistence.SaveState()

}
