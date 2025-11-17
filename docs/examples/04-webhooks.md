## Работа с вебхуками

SDK включает структуру `yoowebhook.WebhookEvent` для обработки событий вебхуков от платежной системы. Она поддерживает следующие события:
- `payment.succeeded` — Успешный платёж
- `payment.waiting_for_capture` — Поступление платежа, который нужно подтвердить
- `payment.canceled` — Отмена платежа или ошибка оплаты
- `refund.succeeded` — Успешный возврат денег покупателю

**Данный пример подходит для платежных решений с HTTP Basic Auth.**

Для работы с вебхуками ознакомьтесь с документацией:
[Входящие уведомления](https://yookassa.ru/developers/using-api/webhooks)

---

### Пример обработки вебхуков

Пример демонстрирует обработку вебхуков с фильтрацией IP-адресов через middleware. Фильтрацию можно также настроить на уровне балансировщика или брандмауэра.

**SDK включает встроенную функцию `IsNotificationIPTrusted()`** для проверки IP-адресов YooKassa. Она автоматически использует актуальный список официальных IP-адресов из [документации YooKassa](https://yookassa.ru/developers/using-api/webhooks#ip).

**Рекомендация:** обрабатывайте вебхуки в фоне, чтобы сервер сразу отвечал `200 OK`. Это позволяет избежать тайм-аутов или повторных отправок вебхука от платежной системы.

**Примечание:** Вам нужно [подтвердить](https://yookassa.ru/developers/using-api/webhooks#using), что вы получили уведомление. Для этого ответьте кодом состояния HTTP `200`. ЮKassa проигнорирует всё, что будет находиться в теле или заголовках ответа. Ответы с любыми другими кодами состояний HTTP будут считаться невалидными, и ЮKassa продолжит доставлять уведомление в течение 24 часов, начиная с момента, когда событие произошло.

#### Пример 1: Использование встроенной проверки IP (рекомендуется)

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
	// Создаем маршрутизатор.
	mux := http.NewServeMux()
	mux.HandleFunc("/webhooks", HandleWebhook)

	// Добавляем middleware для проверки IP YooKassa.
	protectedMux := YooKassaIPFilterMiddleware(mux)

	// Запускаем HTTP-сервер.
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", protectedMux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// HandleWebhook обрабатывает запросы вебхука.
func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Парсим JSON-данные из тела запроса.
	var webhookEvent yoowebhook.WebhookEvent
	err := json.NewDecoder(r.Body).Decode(&webhookEvent)
	if err != nil {
		http.Error(w, "Invalid webhook data", http.StatusBadRequest)
		return
	}

	// Логируем данные.
	log.Printf("Webhook обработан: %+v", webhookEvent)
	log.Printf("Тип Webhook: %+v", webhookEvent.Type)
	log.Printf("Событие: %+v", webhookEvent.Event)
	log.Printf("Данные о платеже: %+v", webhookEvent.Object)

	// Возвращаем HTTP-статус 200 OK.
	w.WriteHeader(http.StatusOK)
}

// Middleware для проверки IP-адресов YooKassa (использует встроенную функцию SDK).
func YooKassaIPFilterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteIP := r.RemoteAddr
		log.Printf("Initial remote IP: %s", remoteIP)

		// Проверяем X-Real-IP, если доступен (для использования за прокси/load balancer).
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			log.Printf("Using X-Real-IP header: %s", realIP)
			remoteIP = realIP
		} else if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			// X-Forwarded-For может содержать несколько IP через запятую, берем первый
			ips := strings.Split(forwardedFor, ",")
			remoteIP = strings.TrimSpace(ips[0])
			log.Printf("Using X-Forwarded-For header: %s", remoteIP)
		}

		// Разделяем адрес на хост и порт.
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

		// Проверяем, является ли IP доверенным адресом YooKassa.
		if !yoowebhook.IsNotificationIPTrusted(host) {
			log.Printf("Rejected webhook from untrusted IP: %s", host)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		log.Printf("Accepted webhook from trusted YooKassa IP: %s", host)
		// Передаем управление дальше.
		next.ServeHTTP(w, r)
	})
}
```

#### Пример 2: Использование пользовательских CIDR диапазонов

Если вам нужно добавить дополнительные IP-адреса (например, для тестирования) или использовать свои диапазоны:

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
	// Получаем официальные диапазоны YooKassa и добавляем тестовые адреса.
	allowedCIDRs := yoowebhook.GetTrustedIPRanges()
	allowedCIDRs = append(allowedCIDRs,
		"127.0.0.1/32", // IPv4 localhost для локального тестирования
		"::1/128",      // IPv6 localhost для локального тестирования
	)

	// Создаем маршрутизатор.
	mux := http.NewServeMux()
	mux.HandleFunc("/webhooks", HandleWebhook)

	// Добавляем middleware для фильтрации IP.
	protectedMux := IPFilterMiddleware(mux, allowedCIDRs)

	// Запускаем HTTP-сервер.
	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", protectedMux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// HandleWebhook обрабатывает запросы вебхука.
func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	var webhookEvent yoowebhook.WebhookEvent
	err := json.NewDecoder(r.Body).Decode(&webhookEvent)
	if err != nil {
		http.Error(w, "Invalid webhook data", http.StatusBadRequest)
		return
	}

	log.Printf("Webhook обработан: %+v", webhookEvent)
	log.Printf("Тип Webhook: %+v", webhookEvent.Type)
	log.Printf("Событие: %+v", webhookEvent.Event)
	log.Printf("Данные о платеже: %+v", webhookEvent.Object)

	w.WriteHeader(http.StatusOK)
}

// Middleware для проверки разрешенных IP-адресов.
func IPFilterMiddleware(next http.Handler, allowedCIDRs []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteIP := r.RemoteAddr

		// Проверяем X-Real-IP, если доступен.
		if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			remoteIP = realIP
		} else if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			ips := strings.Split(forwardedFor, ",")
			remoteIP = strings.TrimSpace(ips[0])
		}

		// Разделяем адрес на хост и порт.
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

		// Проверяем, разрешен ли IP-адрес.
		if !IsIPAllowed(host, allowedCIDRs) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Проверяет, входит ли IP-адрес в разрешенные диапазоны.
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


### Тестирование локально

Пример вебхука с уведомлением `payment.waiting_for_capture`:

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
    "description": "Заказ №72",
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