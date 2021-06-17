package auth

import "time"

// storages stores all available storage engines
var storages []storageInitialization

type storageInitialization struct {
	// StorageType, e.g. redis/fdb/mysql...
	StorageType string
	// This function should return ready to work storage.
	Initialize func(args ...string) (storage, error)
}

// storage defines the methods that the storage engine should have to get / set sessions
type storage interface {
	Get(prefix, key string) ([]byte, error)
	Set(prefix, key string, data []byte, duration time.Duration) error
	Delete(prefix, key string) error
}
