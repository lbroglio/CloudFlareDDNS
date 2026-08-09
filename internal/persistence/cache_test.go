package persistence

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/lbroglio/CloudFlareDDNS/internal/config"
)

func cacheTestSetup(t *testing.T) {
	// Reset the state before each test
	state = nil

	testPeristenceRoot, _ := os.MkdirTemp(os.TempDir(), "test_persistence_root")
	t.Setenv("CLOUDFLAREDDNS_PERSISTENCE_ROOT", testPeristenceRoot)
}

func cacheTestTeardown() {
	// Clean up the temporary persistence root directory after each test
	testPeristenceRoot := config.GetPersistenceRootPath()
	os.RemoveAll(testPeristenceRoot)
}

func makeTestCacheFile() {
	// Create a dummy cache file
	cacheFilePath := config.GetPersistenceRootPath() + "/cache.json"
	os.WriteFile(cacheFilePath, []byte(`{"LastKnownIPsForDNS": {"example.com": "1.2.3.4"}}`), 0644)

}

func getTestState() *State {
	return &State{
		LastKnownIPsForDNS: map[string]string{
			"example.com": "1.2.3.4",
		},
	}
}

func TestCacheFileExists_DoesNotExist(t *testing.T) {
	cacheTestSetup(t)

	result := cacheFileExists()
	if result {
		t.Errorf("Expected cacheFileExists w to return false when the file doesn't exist, got true")
	}
}

func TestCacheFileExists_Exists(t *testing.T) {
	cacheTestSetup(t)
	makeTestCacheFile()

	result := cacheFileExists()
	if !result {
		t.Errorf("Expected cacheFileExists to return true when the file exists, got false")
	}

	cacheTestTeardown()
}

func TestReadCacheFile_Success(t *testing.T) {
	cacheTestSetup(t)
	makeTestCacheFile()

	result, err := readCacheFile()
	if err != nil {
		t.Errorf("Expected readCacheFile to succeed, got error: %v", err)
	}

	expectedIP := "1.2.3.4"
	actualIP := result.LastKnownIPsForDNS["example.com"]
	if actualIP != expectedIP {
		t.Errorf("Expected State to contain IP %s for example.com, got %s", expectedIP, actualIP)
	}

	cacheTestTeardown()
}

func TestReadCacheFile_FileNotFound(t *testing.T) {
	cacheTestSetup(t)
	_, err := readCacheFile()
	if err == nil {
		t.Errorf("Expected readCacheFile to return an error when the file doesn't exist, got nil")
	}

	cacheTestTeardown()
}

func TestReadCacheFile_InvalidJSON(t *testing.T) {
	cacheTestSetup(t)

	// Write the cache file with invalid JSON
	cacheFilePath := config.GetPersistenceRootPath() + "/cache.json"
	os.WriteFile(cacheFilePath, []byte(`{"LastKnownIPsForDNS": {"example.com": "1.2.3.4"}`), 0644)

	_, err := readCacheFile()
	if err == nil {
		t.Errorf("Expected readCacheFile to return an error when the file contains invalid JSON, got nil")
	}

	cacheTestTeardown()
}

func TestGetState_InitializesNewState(t *testing.T) {
	cacheTestSetup(t)

	result, err := GetState()
	if err != nil {
		t.Errorf("Expected GetState to succeed, got error: %v", err)
	}

	if len(result.LastKnownIPsForDNS) != 0 {
		t.Errorf("Expected LastKnownIPsForDNS to be empty, got %v", result.LastKnownIPsForDNS)
	}

	if reflect.DeepEqual(result, *state) == false {
		t.Errorf("Expected GetState returned state to be equal to the global state variable")
	}

	cacheTestTeardown()

}

func TestGetState_ReturnsCachedState(t *testing.T) {
	cacheTestSetup(t)
	state = getTestState()

	expectedState := *state

	result, err := GetState()
	if err != nil {
		t.Errorf("Expected GetState to succeed, got error: %v", err)
	}

	if reflect.DeepEqual(result, expectedState) == false {
		t.Errorf("Expected GetState returned state to be %v, got %v", expectedState, result)
	}

	cacheTestTeardown()

}

func TestGetState_ReturnsStateSavedToFile(t *testing.T) {
	cacheTestSetup(t)
	makeTestCacheFile()
	expectedState := State{
		LastKnownIPsForDNS: map[string]string{
			"example.com": "1.2.3.4",
		},
	}

	result, err := GetState()
	if err != nil {
		t.Errorf("Expected GetState to succeed, got error: %v", err)
	}

	if reflect.DeepEqual(result, expectedState) == false {
		t.Errorf("Expected GetState returned state to be %v, got %v", expectedState, result)
	}

	if reflect.DeepEqual(*state, expectedState) == false {
		t.Errorf("Expected global state variable to be %v, got %v", expectedState, *state)
	}

	cacheTestTeardown()

}

func TestUpdateDNSRecordIP(t *testing.T) {
	cacheTestSetup(t)
	state = getTestState()

	err := UpdateDNSRecordIP("example.com", "4.3.2.1")
	if err != nil {
		t.Errorf("Expected UpdateDNSRecordIP to succeed, got error: %v", err)
	}

	if state.LastKnownIPsForDNS["example.com"] != "4.3.2.1" {
		t.Errorf("Expected LastKnownIPsForDNS for example.com to be updated to 4.3.2.1 , got %s", state.LastKnownIPsForDNS["example.com"])
	}

	cacheTestTeardown()
}

func TestUpdateDNSRecordIP_InitializesState(t *testing.T) {
	cacheTestSetup(t)

	err := UpdateDNSRecordIP("example.com", "1.2.3.4")
	if err != nil {
		t.Errorf("Expected UpdateDNSRecordIP to succeed, got error: %v", err)
	}

	if state.LastKnownIPsForDNS["example.com"] != "1.2.3.4" {
		t.Errorf("Expected LastKnownIPsForDNS for example.com to be updated to 1.2.3.4, got %s", state.LastKnownIPsForDNS["example.com"])
	}

	cacheTestTeardown()
}

func TestUpdateState(t *testing.T) {
	cacheTestSetup(t)
	state = getTestState()

	newState := State{
		LastKnownIPsForDNS: map[string]string{
			"example.com": "4.3.2.1",
		},
	}

	UpdateState(newState)

	if reflect.DeepEqual(*state, newState) == false {
		t.Errorf("Expected global state variable to be %v, got %v", newState, *state)
	}

	cacheTestTeardown()
}

func TestSaveState_NilState(t *testing.T) {
	cacheTestSetup(t)
	state = nil

	err := SaveState()
	if err != nil {
		t.Errorf("Expected SaveState to succeed, got error: %v", err)
	}

	cacheTestTeardown()
}

func TestSaveState_Success(t *testing.T) {
	cacheTestSetup(t)
	state = getTestState()

	err := SaveState()
	if err != nil {
		t.Errorf("Expected SaveState to succeed, got error: %v", err)
	}

	// Read the cache file and verify its contents
	cacheFilePath := config.GetPersistenceRootPath() + "/cache.json"
	bytes, err := os.ReadFile(cacheFilePath)
	if err != nil {
		t.Errorf("Expected to read cache file, got error: %v", err)
	}

	var savedState State
	err = json.Unmarshal(bytes, &savedState)
	if err != nil {
		t.Errorf("Expected to unmarshal cache file, got error: %v", err)
	}

	if reflect.DeepEqual(savedState, *state) == false {
		t.Errorf("Expected saved state to be %v, got %v", *state, savedState)
	}
}
