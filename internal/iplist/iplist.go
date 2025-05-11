package iplist

import (
	"context"
	"net/netip"
)

type IPList interface {
	Add(ctx context.Context, subnet string) error
	Remove(ctx context.Context, subnet string) error
	Contains(ctx context.Context, ip string) (bool, error)
}

func ipInSubnet(ipStr, cidr string) (bool, error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return false, err
	}
	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		return false, err
	}
	return prefix.Contains(ip), nil
}
