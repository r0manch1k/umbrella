# Umbrella Signature Server

Сервис цифровых подписей и лицензирования, обеспечивающий выпуск и проверку лицензий на основе RSA-подписей.

---

## Развертывание проекта

### 1. Копирование конфигурации

```bash
make env
```

### 2. Инициализация проекта (сборка и запуск)

```bash
make init
```

Эта команда:

- соберёт Docker-образы;
- поднимет контейнер с приложением;
- создаст ключи RSA, если они не существуют.

### 3. Управление окружением

```bash
make up # Запустить контейнеры
make down # Остановить контейнеры
make restart # Перезапустить контейнеры
make clean # Удалить контейнеры и образы
make clean-volumes # Удалить контейнеры, образы и тома
```

---

## Разработка и тестирование

### Проверка кода на ошибки

```bash
make lint
```

### Автоматическое исправление стиля

```bash
make fix
```

---

## API методы

### POST /v1/license/issue

Выпускает новую лицензию с цифровой подписью RSA.

- Request:

```json
{
    "user_id": "optional-user-123",
    "duration_hours": 24,
    "hw_fingerprint": "ABCD1234EF567890"
}
```

- Response:

```json
{
    "license": "base64-encoded-license-data",
    "signature": "base64-encoded-signature"
}
```

- Коды ответа:

```text
200 Лицензия успешно сгенерирована
400 Ошибка валидации или данных запроса
500 Ошибка генерации или подписи
```

### POST /v1/license/verify

Проверяет подлинность лицензии и подписи RSA.

- Request:

```json
{
    "license": "base64-encoded-license-data",
    "signature": "base64-encoded-signature"
}
```

- Response:

```json
{
    "valid": true,
    "signature": "base64-encoded-signature"
}
```

```json
{
    "valid": false,
    "reason": "invalid signature"
}
```

- Коды ответа:

```text
Код	Описание
200	Проверка выполнена успешно
400	Ошибка формата запроса
500	Ошибка при проверке подписи
```

## Ключи RSA

- При первом запуске сервис сам создаёт приватный и публичный ключи.
- Публичный ключ выводится в логи при генерации.
- Путь к приватному ключу задаётся через .env:

```text
SIGNATURE_PRIVATE_KEY_PATH=/app/keys/private.pem
SIGNATURE_RSA_KEY_BITS=2048
```
