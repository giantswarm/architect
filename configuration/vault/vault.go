// Package vault provides configuration structures for vault.
package vault

import (
	"net/url"
)

// Vault holds the configuration for the installation's Vault setup
type Vault struct {
	// Address is the URL of the Vault API.
	// e.g: 'https://leaseweb-vault-private.giantswarm.io:8200'
	Address url.URL

	// CaTTL is the TTL for CAs.
	// e.g: `86400h`
	CaTTL string

	// TokenTTL is the TTL for tokens.
	// e.g: '720h'
	TokenTTL string
}
