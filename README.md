[![Golang](https://img.shields.io/badge/Go-v1.24.6-EEEEEE?logo=go&logoColor=white&labelColor=00ADD8)](https://go.dev/)
[![License](https://img.shields.io/pypi/l/yookassa.svg)](LICENSE)

<div align="center">
    <h1 align="center">YooKassa API Golang Client Library
    </h1>
    <h3 align="center">Клиент для работы с платежами по <a href="https://yookassa.ru/developers/api">API ЮKassa</a>
    </h3>
    <p align="center">
        Russian | <a href="README.en.md">English</a> 
    </p>
</div>

### Установка
`go get github.com/sanalrt999/yookassa-sdk-go`

### Начало работы
1. Импортируйте модуль
```golang
import "github.com/sanalrt999/yookassa-sdk-go/yookassa"
```
2. Установите данные для конфигурации
```golang
import "github.com/sanalrt999/yookassa-sdk-go/yookassa"

func main() {
    client := yookassa.NewClient('<Идентификатор магазина>', '<Секретный ключ>')	
}
```
3. Вызовите нужный метод API. [Подробнее в документации к API ЮKassa](https://yookassa.ru/developers/api)

## Примеры использования SDK
#### [Настройки SDK API ЮKassa](https://github.com/sanalrt999/yookassa-sdk-go/blob/main/docs/examples/01-configuration.md)
* [Аутентификация](https://github.com/sanalrt999/yookassa-sdk-go/blob/main/docs/examples/01-configuration.md#Аутентификация)
* [Получение информации о магазине](https://github.com/sanalrt999/yookassa-sdk-go/blob/main/docs/examples/01-configuration.md#Получение-информации-о-магазине)
#### [Работа с платежами](https://github.com/sanalrt999/yookassa-sdk-go/blob/main/docs/examples/02-payments.md)
* [Запрос на создание платежа](https://github.com/sanalrt999/yookassa-sdk-go/blob/main/docs/examples/02-payments.md#Запрос-на-создание-платежа)
* [Запрос на подтверждение платежа](https://github.com/sanalrt999/yookassa-sdk-go/blob/main/docs/examples/02-payments.md#Запрос-на-подтверждение-платежа)
* [Запрос на отмену незавершенного платежа](https://github.com/sanalrt999/yookassa-sdk-go/blob/main/docs/examples/02-payments.md#Запрос-на-отмену-незавершенного-платежа)
* [Получить информацию о платеже](https://github.com/sanalrt999/yookassa-sdk-go/blob/main/docs/examples/02-payments.md#Получить-информацию-о-платеже)
* [Получить список платежей с фильтрацией](https://github.com/sanalrt999/yookassa-sdk-go/blob/main/docs/examples/02-payments.md#Получить-список-платежей-с-фильтрацией)
#### [Работа с возвратами](https://github.com/sanalrt999/yookassa-sdk-go/blob/main/docs/examples/03-refunds.md)
* [Запрос на создание возврата](https://github.com/sanalrt999/yookassa-sdk-go/blob/main/docs/examples/03-refunds.md#Запрос-на-создание-возврата)
* [Получить информацию о возврате](https://github.com/sanalrt999/yookassa-sdk-go/blob/main/docs/examples/03-refunds.md#Получить-информацию-о-возврате)
* [Получить список возвратов с фильтрацией](https://github.com/sanalrt999/yookassa-sdk-go/blob/main/docs/examples/03-refunds.md#Получить-список-возвратов-с-фильтрацией)
#### [Работа с вебхуками](https://github.com/sanalrt999/yookassa-sdk-go/blob/main/docs/examples/04-webhooks.md)
* [Пример обработки вебхуков](https://github.com/sanalrt999/yookassa-sdk-go/blob/main/docs/examples/04-webhooks.md#Пример-обработки-вебхуков)
* [Тестирование локально](https://github.com/sanalrt999/yookassa-sdk-go/blob/main/docs/examples/04-webhooks.md#Тестирование-локально)

## Известные ограничения и планы развития

### Текущие ограничения
- **Низкое покрытие тестами** - требуется добавление unit и интеграционных тестов для всех модулей
- **Отсутствие валидации входных данных** - параметры не проверяются на клиенте перед отправкой в API
- **Поддержка только SBP для выплат** - другие типы выплат пока не реализованы (см. TODO в `payouts.go:64`)

### Рекомендации по безопасности
- **Используйте встроенную проверку IP для вебхуков** - SDK предоставляет функцию `IsNotificationIPTrusted()` для верификации IP-адресов YooKassa
- Не логируйте полные данные о платежах в production-окружении (содержат персональные данные)
- Используйте методы с суффиксом `Ctx` для поддержки таймаутов и отмены операций
- Храните секретные ключи в безопасном хранилище (переменные окружения, secrets manager)
- Настройте IP-фильтрацию на уровне балансировщика или firewall в дополнение к проверке в приложении

### Планы развития
- [ ] Увеличение покрытия тестами до 70%+
- [x] ~~Добавление функции проверки подписи вебхуков~~ - реализована проверка IP-адресов (официальный метод YooKassa)
- [ ] Реализация валидации входных параметров
- [ ] Конфигурируемые таймауты и retry-политики для HTTP клиента
- [ ] Поддержка всех типов выплат
- [ ] Structured logging с поддержкой различных уровней


