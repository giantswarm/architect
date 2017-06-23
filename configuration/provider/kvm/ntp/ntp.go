package ntp

import (
	"net"
)

type NTP struct {
	Servers []net.IP
}
