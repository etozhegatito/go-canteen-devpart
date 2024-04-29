# Используем официальный образ Golang как базовый для сборки приложения
FROM golang:1.21.6 as builder

# Установим рабочую директорию в контейнере
WORKDIR /app

# Скопируем файлы go.mod и go.sum для управления зависимостями
COPY go.mod go.sum ./

# Загрузим зависимости. Это позволит использовать кеш слоя Docker, если зависимости не изменялись
RUN go mod download

# Скопируем исходный код приложения в контейнер
COPY . .

# Соберем наше приложение. Мы не храним отладочную информацию, чтобы образ был меньше
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main ./cmd/server

# Используем небольшой образ Alpine для запуска приложения
FROM alpine:latest  

# Установим CA сертификаты для HTTPS запросов
RUN apk --no-cache add ca-certificates

# Устанавливаем рабочую директорию в контейнере
WORKDIR /root/

# Копируем собранный исполняемый файл из предыдущего шага
COPY --from=builder /app/main .

# Откроем порт 8080 для внешних подключений
EXPOSE 8080

# Запустим приложение
CMD ["./main"]
