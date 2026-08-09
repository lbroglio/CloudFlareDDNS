package persistence

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lbroglio/CloudFlareDDNS/internal/config"
)

// State represents the cached state of the application, including,
// the last known IPs for DNS records.
type State struct {
	LastKnownIPsForDNS map[string]string
}

var state *State

// cacheFileExists checks if the cache file exists in the persistence root directory.
func cacheFileExists() bool {
	filePath := config.GetPersistenceRootPath() + "/cache.json"
	_, err := os.Stat(filePath)

	return !os.IsNotExist(err)
}

// readCacheFile reads the cache file and unmarshals its content into a State struct.
func readCacheFile() (*State, error) {
	filePath := config.GetPersistenceRootPath() + "/cache.json"
	bytes, err := os.ReadFile(filePath)

	if err != nil {
		return nil, fmt.Errorf("failed to read cache file: %v", err)
	}

	var state State
	err = json.Unmarshal(bytes, &state)

	if err != nil {
		return nil, fmt.Errorf("failed to parse cache file: %v", err)
	}

	return &state, nil
}

// GetState returns the cached state of the application.
//
// If the state is not already loaded, it attempts to read it from the cache file or
// initializes a new state if the cache file doesn't exist.
func GetState() (State, error) {
	var errReturn error = nil

	if state == nil {
		if cacheFileExists() {
			cacheContent, err := readCacheFile()

			if err != nil {
				fmt.Printf("Error reading cache file: %v\n", err)
				state = &State{
					LastKnownIPsForDNS: make(map[string]string),
				}
				errReturn = fmt.Errorf("failed to read cache file: %v", err)
			} else {
				state = cacheContent
			}
		} else {
			state = &State{
				LastKnownIPsForDNS: make(map[string]string),
			}
		}
	}

	return *state, errReturn
}

// UpdateDNSRecordIP updates the last known IP for a given DNS record in the cached state.
//
// If the state is not already loaded, it attempts to load it first.
func UpdateDNSRecordIP(recordName string, newIP string) error {
	if state == nil {
		_, err := GetState()
		if err != nil {
			return fmt.Errorf("state failed to initialize: %v", err)
		}
	}

	state.LastKnownIPsForDNS[recordName] = newIP
	return nil
}

// UpdateState updates the cached state of the application with a new State struct.
//
// This can be used to replace a large portion of the state at once, rather than updating individual fields.
func UpdateState(s State) {
	*state = s
}

// SaveState saves the current cached state to the cache file in JSON format.
//
// If the state is nil, it does nothing. If the state is not nil,
// it marshals it into JSON and writes it to the cache file.
func SaveState() error {
	// If the state is nil, there's nothing to save.
	if state == nil {
		return nil
	}

	// Marshal the state into JSON format.
	stateStr, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %v", err)
	}

	// Ensure the persistence root directory exists.
	persistenceRoot := config.GetPersistenceRootPath()
	err = os.MkdirAll(persistenceRoot, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create persistence root directory: %v", err)
	}

	// Write the JSON data to the cache file.
	cacheFilePath := persistenceRoot + "/cache.json"
	err = os.WriteFile(cacheFilePath, stateStr, 0644)
	if err != nil {
		return fmt.Errorf("failed to write state to cache file: %v", err)
	}

	return nil
}
