// Package auth provides configuration structures for authentication/authorization services.
package auth

import (
	"github.com/giantswarm/architect/configuration/auth/vault"
)

// Auth holds configuration for authentication/authorization services.
type Auth struct {
	// Vault holds the configuration for the Vault setup.
	vault.Vault
}
