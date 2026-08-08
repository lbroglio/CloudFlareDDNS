package config

import "os"

func getPersistenceRootPath() string {
	return os.Getenv("CLOUDFLAREDDNS_PERSISTENCE_ROOT")
}
