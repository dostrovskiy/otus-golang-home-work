FROM golang:1.22-alpine

WORKDIR /app

# Устанавливаем psql и goose
RUN apk add --no-cache postgresql15-client
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Копируем миграционные файлы
COPY migrations /migrations

# Entrypoint
COPY build/migrations/docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh

ENTRYPOINT ["sh", "/app/docker-entrypoint.sh"]
