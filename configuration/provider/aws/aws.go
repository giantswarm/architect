// Package aws provides configuration structures for the AWS provider.
package aws

import (
	"github.com/giantswarm/architect/configuration/provider/aws/ec2"
)

// AWS holds configuration for the AWS provider.
type AWS struct {
	// EC2 holds configuration for the EC2 specific settings.
	ec2.EC2
}
