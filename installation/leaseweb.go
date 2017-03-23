package installation

import (
	"net/url"

	"github.com/giantswarm/architect/configuration"
	"github.com/giantswarm/architect/configuration/apiservices"
	"github.com/giantswarm/architect/configuration/guestclusters"
	"github.com/giantswarm/architect/configuration/monitoring"
	"github.com/giantswarm/architect/configuration/vault"
)

var Leaseweb = configuration.Installation{
	V1: configuration.V1{
		APIServices: apiservices.APIServices{
			GSAPI: apiservices.GSAPI{
				Address: url.URL{
					Scheme: "https",
					Host:   "api-g8s.giantswarm.io",
				},
			},
		},
		GuestClusters: guestclusters.GuestClusters{
			APIEndpointBase:  "g8s.fra-1.giantswarm.io",
			HyperkubeVersion: "v1.5.2_coreos.0",
		},
		Monitoring: monitoring.Monitoring{
			Prometheus: monitoring.Prometheus{
				Address: url.URL{
					Scheme: "https",
					Host:   "prometheus-g8s.giantswarm.io",
				},
				RetentionPeriod: "336h",
			},
			Testbot: monitoring.Testbot{
				Interval: "5m",
			},
		},
		Vault: vault.Vault{
			Address: url.URL{
				Scheme: "https",
				Host:   "leaseweb-vault-private.giantswarm.io:8200",
			},
			CaTTL:    "86400h",
			TokenTTL: "720h",
		},
	},
}
