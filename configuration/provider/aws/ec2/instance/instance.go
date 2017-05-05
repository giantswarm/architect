// Package instance provides configuration structures for the EC2 instance
// specific settings.
package instance

import (
	"strings"
)

// kind is a private type to ensure only instance types defined in this package
// can be applied to installation configurations. That prevents other packages
// screwing around with instance type configurations.
type kind string

// The following list of instance types was fetched from the following reference.
//
//     curl -s https://raw.githubusercontent.com/powdahound/ec2instances.info/master/www/instances.json | jq '.[].instance_type'
//
const (
	TypeM1Small    kind = "m1.small"
	TypeM1Medium   kind = "m1.medium"
	TypeM1Large    kind = "m1.large"
	TypeM1XLarge   kind = "m1.xlarge"
	TypeC1Medium   kind = "c1.medium"
	TypeC1XLarge   kind = "c1.xlarge"
	TypeCC28XLarge kind = "cc2.8xlarge"
	TypeCG14XLarge kind = "cg1.4xlarge"
	TypeM2XLarge   kind = "m2.xlarge"
	TypeM22XLarge  kind = "m2.2xlarge"
	TypeM24XLarge  kind = "m2.4xlarge"
	TypeCR18XLarge kind = "cr1.8xlarge"
	TypeI2XLarge   kind = "i2.xlarge"
	TypeI22XLarge  kind = "i2.2xlarge"
	TypeI24XLarge  kind = "i2.4xlarge"
	TypeI28XLarge  kind = "i2.8xlarge"
	TypeHI14XLarge kind = "hi1.4xlarge"
	TypeHS18XLarge kind = "hs1.8xlarge"
	TypeT1Micro    kind = "t1.micro"
	TypeT2Nano     kind = "t2.nano"
	TypeT2Micro    kind = "t2.micro"
	TypeT2Small    kind = "t2.small"
	TypeT2Medium   kind = "t2.medium"
	TypeT2Large    kind = "t2.large"
	TypeT2XLarge   kind = "t2.xlarge"
	TypeT22XLarge  kind = "t2.2xlarge"
	TypeM4Large    kind = "m4.large"
	TypeM4XLarge   kind = "m4.xlarge"
	TypeM42XLarge  kind = "m4.2xlarge"
	TypeM44XLarge  kind = "m4.4xlarge"
	TypeM410XLarge kind = "m4.10xlarge"
	TypeM416XLarge kind = "m4.16xlarge"
	TypeM3Medium   kind = "m3.medium"
	TypeM3Large    kind = "m3.large"
	TypeM3XLarge   kind = "m3.xlarge"
	TypeM32XLarge  kind = "m3.2xlarge"
	TypeC4Large    kind = "c4.large"
	TypeC4XLarge   kind = "c4.xlarge"
	TypeC42XLarge  kind = "c4.2xlarge"
	TypeC44XLarge  kind = "c4.4xlarge"
	TypeC48XLarge  kind = "c4.8xlarge"
	TypeC3Large    kind = "c3.large"
	TypeC3XLarge   kind = "c3.xlarge"
	TypeC32XLarge  kind = "c3.2xlarge"
	TypeC34XLarge  kind = "c3.4xlarge"
	TypeC38XLarge  kind = "c3.8xlarge"
	TypeP2XLarge   kind = "p2.xlarge"
	TypeP28XLarge  kind = "p2.8xlarge"
	TypeP216XLarge kind = "p2.16xlarge"
	TypeG22XLarge  kind = "g2.2xlarge"
	TypeG28XLarge  kind = "g2.8xlarge"
	TypeX116XLarge kind = "x1.16xlarge"
	TypeX132XLarge kind = "x1.32xlarge"
	TypeR4Large    kind = "r4.large"
	TypeR4XLarge   kind = "r4.xlarge"
	TypeR42XLarge  kind = "r4.2xlarge"
	TypeR44XLarge  kind = "r4.4xlarge"
	TypeR48XLarge  kind = "r4.8xlarge"
	TypeR416XLarge kind = "r4.16xlarge"
	TypeR3Large    kind = "r3.large"
	TypeR3XLarge   kind = "r3.xlarge"
	TypeR32XLarge  kind = "r3.2xlarge"
	TypeR34XLarge  kind = "r3.4xlarge"
	TypeR38XLarge  kind = "r3.8xlarge"
	TypeI3Large    kind = "i3.large"
	TypeI3XLarge   kind = "i3.xlarge"
	TypeI32XLarge  kind = "i3.2xlarge"
	TypeI34XLarge  kind = "i3.4xlarge"
	TypeI38XLarge  kind = "i3.8xlarge"
	TypeI316XLarge kind = "i3.16xlarge"
	TypeD2XLarge   kind = "d2.xlarge"
	TypeD22XLarge  kind = "d2.2xlarge"
	TypeD24XLarge  kind = "d2.4xlarge"
	TypeD28XLarge  kind = "d2.8xlarge"
	TypeF12XLarge  kind = "f1.2xlarge"
	TypeF116XLarge kind = "f1.16xlarge"
)

// Instance holds configuration for the EC2 instance specific settings.
type Instance struct {
	// Allowed holds configuration for the allowed instance types on EC2.
	Allowed []kind
	// Available holds configuration for the available instance types on EC2.
	Available []kind
	// Capabilities holds configuration for the instance capabilities on EC2.
	Capabilities map[kind]capabilities
	// Default is the default instance type used when launching guest clusters and
	// not specified otherwhise.
	Default kind
}

// ListToString creates a comma separated list using the provided list of
// instance types.
func ListToString(list []kind) string {
	var stringList []string

	for _, l := range list {
		stringList = append(stringList, string(l))
	}

	return strings.Join(stringList, ",")
}
