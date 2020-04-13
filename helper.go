package socks

import (
	"fmt"
	"net"
)

func parseIP(ip string) (net.IP, error) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return nil, fmt.Errorf("invalid ip address %s", ip)
	}
	for i := 0; i < len(ip); i++ {
		switch ip[i] {
		case '.':
			return parsed.To4(), nil
		case ':':
			return parsed.To16(), nil
		}
	}
	return nil, fmt.Errorf("invalid ip address %s", ip)
}
