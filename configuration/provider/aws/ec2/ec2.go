// Package ec2 provides configuration structures for the EC2 specific settings.
package ec2

import (
	"github.com/giantswarm/architect/configuration/provider/aws/ec2/instance"
)

// EC2 holds configuration for the EC2 specific settings.
type EC2 struct {
	// Instance holds configuration for the instance type settings on EC2.
	instance.Instance
}
