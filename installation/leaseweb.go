package installation

import (
	"net/url"
	"time"

	"github.com/giantswarm/architect/configuration"
	"github.com/giantswarm/architect/configuration/apiservices"
	"github.com/giantswarm/architect/configuration/auth"
	"github.com/giantswarm/architect/configuration/guestclusters"
	"github.com/giantswarm/architect/configuration/monitoring"
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
		Auth: auth.Auth{
			Vault: auth.Vault{
				Address: url.URL{
					Scheme: "https",
					Host:   "leaseweb-vault-private.giantswarm.io:8200",
				},
				CaTTL:    10 * 365 * 24 * time.Hour,
				TokenTTL: 30 * 24 * time.Hour,
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
				RetentionPeriod: 2 * 7 * 24 * time.Hour,
			},
			Testbot: monitoring.Testbot{
				Interval: 5 * time.Minute,
			},
		},
	},
}
