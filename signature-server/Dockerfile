FROM golang:1.25.3-alpine3.22 AS base-app
WORKDIR /app

# Обновляем пакеты alpine
RUN apk update \
    && rm -rf /var/cache/apk/* /tmp/* /usr/share/man

# Устанавливаем зависимости приложения
COPY go.mod go.sum ./
RUN go mod download


FROM base-app AS app-dev

# Устанавливем curl для дальнейший установки зависимостей
RUN apk add --no-cache \
        curl \
    && rm -rf /var/cache/apk/* /tmp/* /usr/share/man

# Устанавливаем air для hotreload
RUN curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.63.0 \
    && air -v

# Устанавливаем golangcilint для проверки качества кода
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.6.1 \
    && golangci-lint --version

# Запускаем приложение в режиме hotreload
COPY docker/app/app.dev.sh /usr/local/bin/app.dev.sh
RUN chmod +x /usr/local/bin/app.dev.sh
ENTRYPOINT ["/usr/local/bin/app.dev.sh"]


FROM base-app AS app-prod-builder

ENV GOOS linux
ENV CGO_ENABLED 0

# Копируем все файлы приложения
COPY . .

# Собираем приложение в бинарный файл
RUN go build -o app


FROM alpine:3.22 AS app-prod

# Устанавливаем зависимости для сертификатов
RUN apk add --no-cache ca-certificates

# Копируем бинарный файл приложения из builder stage и запускаем его
COPY --from=app-prod-builder app .
COPY docker/app/app.prod.sh /usr/local/bin/app.prod.sh
RUN chmod +x /usr/local/bin/app.prod.sh
CMD ["/usr/local/bin/app.prod.sh"]


FROM postgres:17.6-alpine3.22 AS postgres

COPY docker/postgres/init-db.sh /docker-entrypoint-initdb.d/init-db.sh
RUN chmod +x /docker-entrypoint-initdb.d/init-db.sh
