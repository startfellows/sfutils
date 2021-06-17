package auth

import "fmt"

var (
	ErrStorageNotSupported   = fmt.Errorf("unfortunally, goauth not supporting requested storage type at this moment")
	ErrStorageNotInitialized = fmt.Errorf("goauth storage is empty. Please, initialize it first")
	ErrEmptyKey              = fmt.Errorf("goauth session key can't be empty")
	ErrSessionNotFound       = fmt.Errorf("goauth session not found")
)
