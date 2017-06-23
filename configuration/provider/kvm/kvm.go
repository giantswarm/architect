package kvm

import (
	"github.com/giantswarm/architect/configuration/provider/kvm/dns"
	"github.com/giantswarm/architect/configuration/provider/kvm/flannel"
	"github.com/giantswarm/architect/configuration/provider/kvm/ingress"
	"github.com/giantswarm/kubernetesd/flag/service/kvm/ntp"
)

type KVM struct {
	DNS     dns.DNS
	Flannel flannel.Flannel
	Ingress ingress.Ingress
	NTP     ntp.NTP
}
