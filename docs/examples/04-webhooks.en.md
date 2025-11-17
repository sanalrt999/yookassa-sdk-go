## Working with Webhooks

The SDK includes the `yoowebhook.WebhookEvent` structure for handling webhook events from the payment system. It supports the following events:
- `payment.succeeded` — Successful payment
- `payment.waiting_for_capture` — Payment received, awaiting confirmation
- `payment.canceled` — Payment cancellation or payment error
- `refund.succeeded` — Successful refund to the customer

**This example is suitable for payment solutions with HTTP Basic Auth.**

For more details on webhooks, refer to the documentation:
[Incoming Notifications](https://yookassa.ru/developers/using-api/webhooks)

---

### Example Webhook Processing

This example demonstrates processing webhooks with IP address filtering via middleware. Filtering can also be configured at the load balancer or firewall level.

**The SDK includes a built-in `IsNotificationIPTrusted()` function** for verifying YooKassa IP addresses. It automatically uses the current list of official IP addresses from the [YooKassa documentation](https://yookassa.ru/developers/using-api/webhooks#ip).

**Recommendation:** Process webhooks in the background so the server immediately responds with `200 OK`. This helps avoid timeouts or repeated webhook delivery attempts by the payment system.

**Note:** You need to [confirm](https://yookassa.ru/developers/using-api/webhooks#using) that you have received the notification. To do this, respond with the HTTP status code `200`. YooKassa will ignore everything in the response body or headers. Responses with any other HTTP status codes will be considered invalid, and YooKassa will continue delivering the notification for 24 hours starting from when the event occurred.

#### Example 1: Using Built-in IP Verification (Recommended)

```go
package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"

	yoowebhook "github.com/sanalrt999/yookassa-sdk-go/yookassa/webhook"
)

func main() {
	// Create a router.
	mux := http.NewServeMux()
	mux.HandleFunc("/webhooks", HandleWebhook)

	// Add middleware for YooKassa IP verification.
	protectedMux := YooKassaIPFilterMiddleware(mux)

	// Start the HTTP server.
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", protectedMux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// HandleWebhook processes webhook requests.
func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Parse JSON data from the request body.
	var webhookEvent yoowebhook.WebhookEvent
	err := json.NewDecoder(r.Body).Decode(&webhookEvent)
	if err != nil {
		http.Error(w, "Invalid webhook data", http.StatusBadRequest)
		return
	}

	// Log data.
	log.Printf("Webhook processed: %+v", webhookEvent)
	log.Printf("Webhook Type: %+v", webhookEvent.Type)
	log.Printf("Event: %+v", webhookEvent.Event)
	log.Printf("Payment Data: %+v", webhookEvent.Object)

	// Return HTTP status 200 OK.
	w.WriteHeader(http.StatusOK)
}

// Middleware to verify YooKassa IP addresses (uses the SDK's built-in function).
func YooKassaIPFilterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteIP := r.RemoteAddr
		log.Printf("Initial remote IP: %s", remoteIP)

		// Check X-Real-IP header if available (for use behind proxy/load balancer).
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			log.Printf("Using X-Real-IP header: %s", realIP)
			remoteIP = realIP
		} else if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			// X-Forwarded-For may contain multiple IPs separated by commas, take the first one
			ips := strings.Split(forwardedFor, ",")
			remoteIP = strings.TrimSpace(ips[0])
			log.Printf("Using X-Forwarded-For header: %s", remoteIP)
		}

		// Split the address into host and port.
		var host string
		if strings.Contains(remoteIP, ":") {
			var err error
			host, _, err = net.SplitHostPort(remoteIP)
			if err != nil {
				log.Printf("Failed to parse IP address: %v", err)
				http.Error(w, "Invalid remote IP address", http.StatusBadRequest)
				return
			}
		} else {
			host = remoteIP
		}

		// Check if the IP is a trusted YooKassa address.
		if !yoowebhook.IsNotificationIPTrusted(host) {
			log.Printf("Rejected webhook from untrusted IP: %s", host)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		log.Printf("Accepted webhook from trusted YooKassa IP: %s", host)
		// Pass control to the next handler.
		next.ServeHTTP(w, r)
	})
}
```

#### Example 2: Using Custom CIDR Ranges

If you need to add additional IP addresses (e.g., for testing) or use your own ranges:

```go
package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"

	yoowebhook "github.com/sanalrt999/yookassa-sdk-go/yookassa/webhook"
)

func main() {
	// Get official YooKassa ranges and add test addresses.
	allowedCIDRs := yoowebhook.GetTrustedIPRanges()
	allowedCIDRs = append(allowedCIDRs,
		"127.0.0.1/32", // IPv4 localhost for local testing
		"::1/128",      // IPv6 localhost for local testing
	)

	// Create a router.
	mux := http.NewServeMux()
	mux.HandleFunc("/webhooks", HandleWebhook)

	// Add middleware for IP filtering.
	protectedMux := IPFilterMiddleware(mux, allowedCIDRs)

	// Start the HTTP server.
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", protectedMux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// HandleWebhook processes webhook requests.
func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var webhookEvent yoowebhook.WebhookEvent
	err := json.NewDecoder(r.Body).Decode(&webhookEvent)
	if err != nil {
		http.Error(w, "Invalid webhook data", http.StatusBadRequest)
		return
	}

	log.Printf("Webhook processed: %+v", webhookEvent)
	log.Printf("Webhook Type: %+v", webhookEvent.Type)
	log.Printf("Event: %+v", webhookEvent.Event)
	log.Printf("Payment Data: %+v", webhookEvent.Object)

	w.WriteHeader(http.StatusOK)
}

// Middleware to check allowed IP addresses.
func IPFilterMiddleware(next http.Handler, allowedCIDRs []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteIP := r.RemoteAddr

		// Check X-Real-IP header if available.
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			remoteIP = realIP
		} else if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			ips := strings.Split(forwardedFor, ",")
			remoteIP = strings.TrimSpace(ips[0])
		}

		// Split the address into host and port.
		var host string
		if strings.Contains(remoteIP, ":") {
			var err error
			host, _, err = net.SplitHostPort(remoteIP)
			if err != nil {
				http.Error(w, "Invalid remote IP address", http.StatusBadRequest)
				return
			}
		} else {
			host = remoteIP
		}

		// Check if the IP address is allowed.
		if !IsIPAllowed(host, allowedCIDRs) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Check if the IP address is within the allowed ranges.
func IsIPAllowed(ip string, allowedCIDRs []string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, cidr := range allowedCIDRs {
		_, allowedNet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if allowedNet.Contains(parsedIP) {
			return true
		}
	}
	return false
}
```

### Local Testing

Example webhook notification with `payment.waiting_for_capture`:

```bash
curl -i --location 'http://localhost:8080/webhooks' \
--header 'Content-Type: application/json' \
--data '{
  "type": "notification",
  "event": "payment.waiting_for_capture",
  "object": {
    "id": "22d6d597-000f-5000-9000-145f6df21d6f",
    "status": "waiting_for_capture",
    "paid": true,
    "amount": {
      "value": "2.00",
      "currency": "RUB"
    },
    "authorization_details": {
      "rrn": "10000000000",
      "auth_code": "000000",
      "three_d_secure": {
        "applied": true
      }
    },
    "created_at": "2018-07-10T14:27:54.691Z",
    "description": "Order #72",
    "expires_at": "2018-07-17T14:28:32.484Z",
    "metadata": {},
    "payment_method": {
      "type": "bank_card",
      "id": "22d6d597-000f-5000-9000-145f6df21d6f",
      "saved": false,
      "card": {
        "first6": "555555",
        "last4": "4444",
        "expiry_month": "07",
        "expiry_year": "2021",
        "card_type": "MasterCard",
        "issuer_country": "RU",
        "issuer_name": "Sberbank"
      },
      "title": "Bank card *4444"
    },
    "refundable": false,
    "test": false
  }
}
'
```
