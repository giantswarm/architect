package installation

import (
	"net/url"
	"time"

	"github.com/giantswarm/architect/configuration"
	"github.com/giantswarm/architect/configuration/auth"
	"github.com/giantswarm/architect/configuration/cluster"
	"github.com/giantswarm/architect/configuration/giantswarm"
	"github.com/giantswarm/architect/configuration/monitoring"
)

var Leaseweb = configuration.Installation{
	V1: configuration.V1{
		Auth: auth.Auth{
			Vault: auth.Vault{
				Address: url.URL{
					Scheme: "https",
					Host:   "leaseweb-vault-private.giantswarm.io:8200",
				},
				CA: auth.CA{
					TTL: 10 * 365 * 24 * time.Hour,
				},
				Token: auth.Token{
					TTL: 30 * 24 * time.Hour,
				},
			},
		},

		GiantSwarm: giantswarm.GiantSwarm{
			API: giantswarm.API{
				Address: url.URL{
					Scheme: "https",
					Host:   "api-g8s.giantswarm.io",
				},
			},
		},

		Guest: cluster.Guest{
			Hyperkube: cluster.Hyperkube{
				Version: "v1.5.2_coreos.0",
			},
			Kubernetes: cluster.Kubernetes{
				API: cluster.API{
					EndpointBase: "g8s.fra-1.giantswarm.io",
				},
			},
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
