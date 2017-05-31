package route53

import (
	"github.com/giantswarm/architect/configuration/provider/aws/route53/hostedzones"
)

// Route53 holds configuration for the Route53 specific settings.
type Route53 struct {
	// HostedZones holds the Hosted Zone IDs for creating guest cluster
	// recordsets.
	HostedZones hostedzones.HostedZones
}
