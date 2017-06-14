package kvm

import "github.com/giantswarm/architect/configuration/provider/kvm/flannel"

type KVM struct {
	Flannel flannel.Flannel
}
