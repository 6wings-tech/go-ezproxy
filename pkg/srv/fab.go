package srv

import "net/netip"

func New(addr netip.AddrPort) *server {
	return &server{
		addr: addr.String(),
	}
}
