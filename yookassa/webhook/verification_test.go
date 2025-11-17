package yoowebhook

import (
	"testing"
)

func TestIsNotificationIPTrusted(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		// Valid YooKassa IPv4 addresses
		{
			name:     "Valid IP from 185.71.76.0/27 range",
			ip:       "185.71.76.5",
			expected: true,
		},
		{
			name:     "First IP in 185.71.76.0/27 range",
			ip:       "185.71.76.0",
			expected: true,
		},
		{
			name:     "Last IP in 185.71.76.0/27 range",
			ip:       "185.71.76.31",
			expected: true,
		},
		{
			name:     "Valid IP from 185.71.77.0/27 range",
			ip:       "185.71.77.15",
			expected: true,
		},
		{
			name:     "Valid IP from 77.75.153.0/25 range",
			ip:       "77.75.153.100",
			expected: true,
		},
		{
			name:     "First IP in 77.75.153.0/25 range",
			ip:       "77.75.153.0",
			expected: true,
		},
		{
			name:     "Last IP in 77.75.153.0/25 range",
			ip:       "77.75.153.127",
			expected: true,
		},
		{
			name:     "Valid IP from 77.75.154.128/25 range",
			ip:       "77.75.154.200",
			expected: true,
		},
		{
			name:     "First IP in 77.75.154.128/25 range",
			ip:       "77.75.154.128",
			expected: true,
		},
		{
			name:     "Last IP in 77.75.154.128/25 range",
			ip:       "77.75.154.255",
			expected: true,
		},
		{
			name:     "Individual IP 77.75.156.11",
			ip:       "77.75.156.11",
			expected: true,
		},
		{
			name:     "Individual IP 77.75.156.35",
			ip:       "77.75.156.35",
			expected: true,
		},
		// Valid YooKassa IPv6 addresses
		{
			name:     "Valid IPv6 from 2a02:5180::/32 range",
			ip:       "2a02:5180::1",
			expected: true,
		},
		{
			name:     "Valid IPv6 from 2a02:5180::/32 range (full notation)",
			ip:       "2a02:5180:0:0:0:0:0:1",
			expected: true,
		},
		{
			name:     "First IPv6 in 2a02:5180::/32 range",
			ip:       "2a02:5180::",
			expected: true,
		},
		{
			name:     "IPv6 within 2a02:5180::/32 range",
			ip:       "2a02:5180:1234:5678:9abc:def0:1234:5678",
			expected: true,
		},
		// Invalid/untrusted IP addresses
		{
			name:     "Random public IPv4",
			ip:       "8.8.8.8",
			expected: false,
		},
		{
			name:     "IP just outside 185.71.76.0/27 range",
			ip:       "185.71.76.32",
			expected: false,
		},
		{
			name:     "IP just before 185.71.76.0/27 range",
			ip:       "185.71.75.255",
			expected: false,
		},
		{
			name:     "IP in same /24 but outside 77.75.153.0/25",
			ip:       "77.75.153.128",
			expected: false,
		},
		{
			name:     "IP just before 77.75.154.128/25 range",
			ip:       "77.75.154.127",
			expected: false,
		},
		{
			name:     "Individual IP close to 77.75.156.11 but not exact",
			ip:       "77.75.156.10",
			expected: false,
		},
		{
			name:     "Individual IP close to 77.75.156.35 but not exact",
			ip:       "77.75.156.36",
			expected: false,
		},
		{
			name:     "Localhost IPv4",
			ip:       "127.0.0.1",
			expected: false,
		},
		{
			name:     "Localhost IPv6",
			ip:       "::1",
			expected: false,
		},
		{
			name:     "Private IPv4 192.168.x.x",
			ip:       "192.168.1.1",
			expected: false,
		},
		{
			name:     "Private IPv4 10.x.x.x",
			ip:       "10.0.0.1",
			expected: false,
		},
		{
			name:     "IPv6 outside 2a02:5180::/32 range",
			ip:       "2a02:5181::1",
			expected: false,
		},
		{
			name:     "Random IPv6",
			ip:       "2001:4860:4860::8888",
			expected: false,
		},
		// Invalid IP formats
		{
			name:     "Invalid IP format",
			ip:       "256.256.256.256",
			expected: false,
		},
		{
			name:     "Empty string",
			ip:       "",
			expected: false,
		},
		{
			name:     "Invalid characters",
			ip:       "not-an-ip",
			expected: false,
		},
		{
			name:     "IP with port",
			ip:       "185.71.76.5:8080",
			expected: false,
		},
		{
			name:     "IPv6 with port bracket notation",
			ip:       "[2a02:5180::1]:8080",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNotificationIPTrusted(tt.ip)
			if result != tt.expected {
				t.Errorf("IsNotificationIPTrusted(%q) = %v, expected %v", tt.ip, result, tt.expected)
			}
		})
	}
}

func TestGetTrustedIPRanges(t *testing.T) {
	ranges := GetTrustedIPRanges()

	// Check that we got the expected number of ranges
	expectedCount := 7 // 4 IPv4 CIDR ranges + 2 individual IPs + 1 IPv6 CIDR range
	if len(ranges) != expectedCount {
		t.Errorf("GetTrustedIPRanges() returned %d ranges, expected %d", len(ranges), expectedCount)
	}

	// Verify some expected ranges are present
	expectedRanges := map[string]bool{
		"185.71.76.0/27":    false,
		"185.71.77.0/27":    false,
		"77.75.153.0/25":    false,
		"77.75.154.128/25":  false,
		"77.75.156.11/32":   false,
		"77.75.156.35/32":   false,
		"2a02:5180::/32":    false,
	}

	for _, cidr := range ranges {
		if _, exists := expectedRanges[cidr]; exists {
			expectedRanges[cidr] = true
		}
	}

	for cidr, found := range expectedRanges {
		if !found {
			t.Errorf("Expected CIDR %q not found in GetTrustedIPRanges() result", cidr)
		}
	}

	// Verify that modifying the returned slice doesn't affect the original
	originalLen := len(ranges)
	ranges[0] = "0.0.0.0/0"
	newRanges := GetTrustedIPRanges()
	if len(newRanges) != originalLen {
		t.Error("Modifying returned slice affected the original data")
	}
	if newRanges[0] == "0.0.0.0/0" {
		t.Error("Modifying returned slice affected subsequent calls")
	}
}

// Benchmark for IsNotificationIPTrusted with valid IP
func BenchmarkIsNotificationIPTrusted_Valid(b *testing.B) {
	ip := "185.71.76.5"
	for i := 0; i < b.N; i++ {
		IsNotificationIPTrusted(ip)
	}
}

// Benchmark for IsNotificationIPTrusted with invalid IP
func BenchmarkIsNotificationIPTrusted_Invalid(b *testing.B) {
	ip := "8.8.8.8"
	for i := 0; i < b.N; i++ {
		IsNotificationIPTrusted(ip)
	}
}

// Benchmark for IsNotificationIPTrusted with IPv6
func BenchmarkIsNotificationIPTrusted_IPv6(b *testing.B) {
	ip := "2a02:5180::1"
	for i := 0; i < b.N; i++ {
		IsNotificationIPTrusted(ip)
	}
}
