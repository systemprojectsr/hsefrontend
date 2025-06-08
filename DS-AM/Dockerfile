FROM golang:1.24-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости для компиляции
RUN apk add --no-cache ca-certificates git

# Копируем go.mod и go.sum для кеширования слоя с зависимостями
COPY go.mod ./
RUN go mod download

# Копируем весь код проекта
COPY . .

# Компилируем сервис
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o photoservice ./cmd/photoservice/main.go

# Многоэтапная сборка для минимизации размера финального образа
FROM alpine:3.19

# Обновление сертификатов CA
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Копируем скомпилированный бинарник и конфигурацию
COPY --from=builder /app/photoservice /app/
COPY --from=builder /app/configs/photoservice.json /app/configs/

# Открываем порт для API
EXPOSE 8081

# Запускаем сервис
CMD ["/app/photoservice"] 