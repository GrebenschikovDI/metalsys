package subnet

import "net"

func IsTrusted(ip, trustedSubnet string) bool {
	_, trustedIPNet, err := net.ParseCIDR(trustedSubnet)
	if err != nil {
		return false
	}

	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return false
	}

	return trustedIPNet.Contains(clientIP)
}
