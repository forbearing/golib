package controller

import (
	"net"
	"strings"
)

// IPv6ToIPv4 converts IPv6 to IPv4 if possible
func IPv6ToIPv4(ipStr string) string {
	// If its ipv4, return.
	if net.ParseIP(ipStr).To4() != nil {
		return ipStr
	}

	// handle IPv6 localhost
	if strings.HasPrefix(ipStr, "::") {
		return "127.0.0.1"
	}

	// handle IPv4-mapped IPv6 addresses
	// eg ::ffff:192.0.2.128 或 ::ffff:c000:280
	if strings.Contains(ipStr, "::ffff:") {
		split := strings.Split(ipStr, "::ffff:")
		if len(split) == 2 {
			if ip := net.ParseIP(split[1]).To4(); ip != nil {
				return ip.String()
			}
		}
	}

	// handle embedded IPv4 address
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ipStr
	}

	ip4 := ip.To4()
	if ip4 != nil {
		return ip4.String()
	}

	return ipStr
}
