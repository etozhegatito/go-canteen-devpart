# Используйте официальный образ Go как базовый
FROM golang:1.21.6 as builder


# Установите рабочий каталог в контейнере
WORKDIR /app


# Копируйте go модуль и сумму файлов
COPY . .
# Копируйте все файлы проекта в контейнер
COPY . .

# Соберите приложение
RUN go build -o app ./cmd/server/main.go

EXPOSE 8080

CMD ["./app"]

