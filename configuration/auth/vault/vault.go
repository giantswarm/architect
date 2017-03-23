// Package vault provides configuration structures for Vault.
package vault

import (
	"net/url"
	"time"
)

// Vault holds the configuration for the installation's Vault setup
type Vault struct {
	// Address is the URL of the Vault API.
	// e.g: 'https://leaseweb-vault-private.giantswarm.io:8200'
	Address url.URL

	// CA is the configuration for Vault CAs.
	CA

	// Token is the configuration for Vault tokens.
	Token
}

// CA holds configuration for a certificate authority.
type CA struct {
	// TTL is the TTL for the CA.
	TTL time.Duration
}

// Token holds configuration for a Vault token.
type Token struct {
	// TTL is the TTL for the Token.
	TTL time.Duration
}
