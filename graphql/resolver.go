package graphql

import (
	"crypto/ecdsa"
	"github.com/pixlcrashr/squadmortar.xyz-sessions-server/session/storage"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	EcdsaKey       *ecdsa.PrivateKey
	SessionStorage storage.Storage
}
