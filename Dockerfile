FROM golang:1.21.6 as builder


#рабочий каталог в контейнере
WORKDIR /app


#Копипаста go модуль
COPY . .

# для сборки приложение
RUN go build -o app ./cmd/server/main.go

#Занимаем этот порт в контейнере
EXPOSE 8070

#Запускаем
CMD ["./app"]

