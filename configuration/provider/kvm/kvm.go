package kvm

import (
	"github.com/giantswarm/architect/configuration/provider/kvm/flannel"
	"github.com/giantswarm/architect/configuration/provider/kvm/ingress"
)

type KVM struct {
	Flannel flannel.Flannel
	Ingress ingress.Ingress
}
