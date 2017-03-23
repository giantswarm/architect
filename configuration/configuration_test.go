package configuration

import (
	"net/url"
	"testing"

	"github.com/giantswarm/architect/configuration/auth"
	"github.com/giantswarm/architect/configuration/auth/vault"
)

// TestAccessingVaultAddress tests that the configuration
// path works as expected.
func TestAccessingVaultAddress(t *testing.T) {
	config := Installation{
		V1{
			Auth: auth.Auth{
				Vault: vault.Vault{
					Address: url.URL{
						Host: "foo.bar.com",
					},
				},
			},
		},
	}

	if config.V1.Auth.Vault.Address.Host != "foo.bar.com" {
		t.Fatalf("could not access vault address")
	}
}
