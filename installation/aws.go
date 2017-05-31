package installation

import (
	"net/url"
	"time"

	"github.com/giantswarm/architect/configuration"
	"github.com/giantswarm/architect/configuration/auth"
	"github.com/giantswarm/architect/configuration/auth/vault"
	"github.com/giantswarm/architect/configuration/giantswarm"
	"github.com/giantswarm/architect/configuration/giantswarm/api"
	"github.com/giantswarm/architect/configuration/giantswarm/desmotes"
	"github.com/giantswarm/architect/configuration/giantswarm/happa"
	"github.com/giantswarm/architect/configuration/giantswarm/passage"
	"github.com/giantswarm/architect/configuration/guest"
	"github.com/giantswarm/architect/configuration/guest/hyperkube"
	"github.com/giantswarm/architect/configuration/guest/kubectl"
	"github.com/giantswarm/architect/configuration/guest/kubernetes"
	"github.com/giantswarm/architect/configuration/monitoring"
	"github.com/giantswarm/architect/configuration/monitoring/prometheus"
	"github.com/giantswarm/architect/configuration/monitoring/testbot"
	"github.com/giantswarm/architect/configuration/provider"
	"github.com/giantswarm/architect/configuration/provider/aws"
	"github.com/giantswarm/architect/configuration/provider/aws/ec2"
	"github.com/giantswarm/architect/configuration/provider/aws/ec2/instance"
)

var AWS = configuration.Installation{
	V1: configuration.V1{
		Auth: auth.Auth{
			Vault: vault.Vault{
				Address: url.URL{
					Scheme: "https",
					Host:   "vault.eu-west-1.aws.private.giantswarm.io",
				},
				CA: vault.CA{
					TTL: 10 * 365 * 24 * time.Hour,
				},
				Certificate: vault.Certificate{
					TTL: 6 * 30 * 24 * time.Hour,
				},
				Token: vault.Token{
					TTL: 6 * 30 * 24 * time.Hour,
				},
			},
		},

		GiantSwarm: giantswarm.GiantSwarm{
			API: api.API{
				Address: url.URL{
					Scheme: "https",
					Host:   "api-aws.giantswarm.io",
				},
			},
			Passage: passage.Passage{
				Address: url.URL{
					Scheme: "https",
					Host:   "passage-aws.giantswarm.io",
				},
			},
			Desmotes: desmotes.Desmotes{
				Address: url.URL{
					Scheme: "https",
					Host:   "desmotes-aws.giantswarm.io",
				},
			},
			Happa: happa.Happa{
				Address: url.URL{
					Scheme: "https",
					Host:   "happa-aws.giantswarm.io",
				},
			},
		},

		Guest: guest.Guest{
			Hyperkube: hyperkube.Hyperkube{
				Version: hyperkube.Version,
			},
			Kubectl: kubectl.Kubectl{
				Version: kubectl.Version,
			},
			Kubernetes: kubernetes.Kubernetes{
				API: kubernetes.API{
					EndpointBase: "g8s.eu-west-1.adidas.aws.giantswarm.io",
				},
			},
		},

		Monitoring: monitoring.Monitoring{
			Prometheus: prometheus.Prometheus{
				Address: url.URL{
					Scheme: "https",
					Host:   "prometheus-aws.giantswarm.io",
				},
				RetentionPeriod: 2 * 7 * 24 * time.Hour,
			},
			Testbot: testbot.Testbot{
				Interval: 30 * time.Minute,
			},
		},

		Provider: provider.Provider{
			AWS: aws.AWS{
				EC2: ec2.EC2{
					Instance: instance.Instance{
						Allowed: instance.Allowed(
							instance.TypeM3Large,
							instance.TypeM3XLarge,
							instance.TypeM32XLarge,
							instance.TypeR3Large,
							instance.TypeR3XLarge,
							instance.TypeR32XLarge,
							instance.TypeR34XLarge,
							instance.TypeR38XLarge,
						),
						Available:    instance.Available(),
						Capabilities: instance.Capabilities(),
						Default:      instance.Default,
					},
				},
			},
			Kind: provider.AWS,
		},
	},
}
