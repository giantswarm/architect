// Package auth provides configuration structures for authentication/authorization services.
package auth

import (
	"net/url"
	"time"
)

// Auth holds configuration for authentication/authorization services.
type Auth struct {
	Vault
}

// Vault holds the configuration for the installation's Vault setup
type Vault struct {
	// Address is the URL of the Vault API.
	// e.g: 'https://leaseweb-vault-private.giantswarm.io:8200'
	Address url.URL

	// CaTTL is the TTL for CAs.
	// e.g: `86400h`
	CaTTL time.Duration

	// TokenTTL is the TTL for tokens.
	// e.g: '720h'
	TokenTTL time.Duration
}
