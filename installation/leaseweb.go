package installation

import (
	"github.com/giantswarm/architect/configuration"
)

var Leaseweb = configuration.Installation{
	V1: configuration.V1{
		APIServices: configuration.APIServices{
			GSAPI: configuration.GSAPI{
				Address: "https://api-g8s.giantswarm.io",
			},
		},
		GuestClusters: configuration.GuestClusters{
			APIEndpointBase:  "g8s.fra-1.giantswarm.io",
			HyperkubeVersion: "v1.5.2_coreos.0",
		},
		Monitoring: configuration.Monitoring{
			Prometheus: configuration.Prometheus{
				Address:         "https://prometheus-g8s.giantswarm.io",
				RetentionPeriod: "336h",
			},
			Testbot: configuration.Testbot{
				Interval: "5m",
			},
		},
		Vault: configuration.Vault{
			Address:  "https://leaseweb-vault-private.giantswarm.io:8200",
			CaTTL:    "86400h",
			TokenTTL: "720h",
		},
	},
}
