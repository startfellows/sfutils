package auth

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// Auth main structure providing access to receiving and storing sessions
type Auth struct {
	storage storage
	prefix  string
	seed    *rand.Rand
}

// generateKey returns a pseudo-random key for a new session
func (a *Auth) generateKey() string {
	randomBytes := make([]byte, 16)
	_, err := a.seed.Read(randomBytes)
	if err != nil {
		panic(fmt.Errorf("goauth fatal error during rand.Read: %s", err.Error()))
	}

	return hex.EncodeToString(randomBytes)
}

// Get receives data via key and unmarshal them to the interface i
func (a *Auth) Get(key string, i interface{}) error {
	if a.storage == nil {
		return ErrStorageNotInitialized
	}
	if len(key) == 0 {
		return ErrEmptyKey
	}

	data, err := a.storage.Get(a.prefix, key)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &i)
	if err != nil {
		return err
	}

	return nil
}

// Add adds a new session and returns a new key on success
func (a *Auth) Add(value interface{}, duration time.Duration) (string, error) {
	newKey := a.generateKey()

	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	err = a.storage.Set(a.prefix, newKey, data, duration)
	if err != nil {
		return "", err
	}

	return newKey, nil
}

// Delete just removes session by key
func (a *Auth) Delete(key string) error {
	err := a.storage.Delete(a.prefix, key)
	if err != nil {
		return err
	}

	return nil
}

// New will create new Auth structure.
// storageType currently supports only redis and fdb.
// prefix is used to separate different types of sessions.
// For example, you want to have sessions for moderators and users. And accordingly, you need two prefixes, "moder" and "user".
// This is also an additional measure against possible attacks and collisions.
// storageObject is object, that will be used for data marshaling/unmarshalling in storage.
// args is storage specific:
// redis: 2 strings. first is connection string, second is password (if exist)
// fdb: nothing, instead fdb using environment variable FDB_CLUSTER_FILE
func New(storageType, prefix string, args ...string) (*Auth, error) {
	newAuth := &Auth{
		prefix: prefix,
		seed:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	for _, s := range storages {
		if s.StorageType == storageType {
			var err error
			newAuth.storage, err = s.Initialize(args...)
			if err != nil {
				return nil, err
			}

			break
		}
	}

	if newAuth.storage == nil {
		return nil, ErrStorageNotSupported
	}

	return newAuth, nil
}
