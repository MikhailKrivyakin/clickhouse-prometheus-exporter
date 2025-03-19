# Используем официальный образ Go
FROM golang:1.19-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o clickhouse-prometheus-exporter ./cmd/exporter/main.go

# Используем минимальный образ для финального контейнера
FROM alpine:latest

# Копируем собранное приложение
COPY --from=builder /app/clickhouse-prometheus-exporter /clickhouse-prometheus-exporter
COPY config/queries.yaml /config/queries.yaml

# Указываем порт, который будет использоваться
EXPOSE 8080

# Команда для запуска приложения
CMD ["/clickhouse-prometheus-exporter"]