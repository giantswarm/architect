// Package aws provides configuration structures for the AWS provider.
package aws

import (
	"github.com/giantswarm/architect/configuration/provider/aws/ec2"
	"github.com/giantswarm/architect/configuration/provider/aws/route53"
)

// AWS holds configuration for the AWS provider.
type AWS struct {
	// EC2 holds configuration for the EC2 specific settings.
	ec2.EC2
	// Route53 holds configuration for the Route53 specific settings.
	route53.Route53
}
