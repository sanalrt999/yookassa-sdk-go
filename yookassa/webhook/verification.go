package yoowebhook

import (
	"fmt"
	"net"
)

// YooKassa official IP addresses and CIDR ranges from which webhooks are sent.
// Documentation: https://yookassa.ru/developers/using-api/webhooks#ip
var trustedCIDRs = []string{
	// IPv4 CIDR ranges
	"185.71.76.0/27",
	"185.71.77.0/27",
	"77.75.153.0/25",
	"77.75.154.128/25",
	// Individual IPv4 addresses
	"77.75.156.11/32",
	"77.75.156.35/32",
	// IPv6 CIDR range
	"2a02:5180::/32",
}

// trustedNetworks caches parsed CIDR networks for performance
var trustedNetworks []*net.IPNet

func init() {
	trustedNetworks = make([]*net.IPNet, 0, len(trustedCIDRs))
	for _, cidr := range trustedCIDRs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			// This should never happen with hardcoded valid CIDRs
			panic(fmt.Sprintf("invalid CIDR in trustedCIDRs: %s: %v", cidr, err))
		}
		trustedNetworks = append(trustedNetworks, network)
	}
}

// IsNotificationIPTrusted checks whether an IP address is among the IP addresses
// from which YooKassa sends webhook notifications.
//
// Parameters:
//   - ip: string containing an IPv4 or IPv6 address (e.g., "185.71.76.5" or "2a02:5180::1")
//
// Returns:
//   - true if the IP address belongs to YooKassa's trusted IP ranges
//   - false otherwise
//
// Example:
//
//	if !yoowebhook.IsNotificationIPTrusted("185.71.76.5") {
//	    http.Error(w, "Forbidden", http.StatusForbidden)
//	    return
//	}
func IsNotificationIPTrusted(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, network := range trustedNetworks {
		if network.Contains(parsedIP) {
			return true
		}
	}
	return false
}

// GetTrustedIPRanges returns a copy of the list of trusted CIDR ranges
// from which YooKassa sends webhook notifications.
//
// This function can be useful for configuring firewalls or load balancers.
//
// Returns:
//   - A slice of strings containing CIDR notation (e.g., ["185.71.76.0/27", "185.71.77.0/27", ...])
func GetTrustedIPRanges() []string {
	// Return a copy to prevent modification
	result := make([]string, len(trustedCIDRs))
	copy(result, trustedCIDRs)
	return result
}
