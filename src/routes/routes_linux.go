//go:build linux

package routes

import (
	"fmt"

	"github.com/gologme/log"
	"github.com/vishvananda/netlink"

	"github.com/yggdrasil-network/yggdrasil-go/src/tun"
)

func SetRoutes(tun *tun.TunAdapter, log *log.Logger, cidrs []string, src string) error {
	srcaddr, err := netlink.ParseAddr(src)
	if err != nil {
		return fmt.Errorf("couldn't parse source CIDR %q: %w", src, err)
	}
	nlintf, err := netlink.LinkByName(tun.Name())
	if err != nil {
		return fmt.Errorf("failed to find link by name: %w", err)
	}
	for _, cidr := range cidrs {
		nladdr, err := netlink.ParseAddr(cidr)
		if err != nil {
			return fmt.Errorf("couldn't parse CIDR %q: %w", cidr, err)
		}
		if err := netlink.RouteAdd(&netlink.Route{
			Src:       srcaddr.IPNet.IP,
			Dst:       nladdr.IPNet,
			LinkIndex: nlintf.Attrs().Index,
			Scope:     netlink.SCOPE_LINK,
		}); err != nil {
			log.Warnln("Failed to add route", cidr, "to routing table:", err)
		}
	}
	return nil
}

func AddIP(iface string, addr string) error {
	lo, err := netlink.LinkByName(iface)
	if err != nil {
		return fmt.Errorf("couldn't find interface %q: %w", iface, err)
	}
	loaddr, err := netlink.ParseAddr(addr)
	if err != nil {
		return fmt.Errorf("couldn't parse source CIDR %q: %w", addr, err)
	}
	netlink.AddrAdd(lo, loaddr)
	return nil
}