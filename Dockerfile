#Dockerfile.multistage
FROM golang:alpine AS builder

## Build
LABEL stage=gobuilder

#cgo отключен по умолчанию
ENV CGO_ENABLED 0

ENV GOOS linux

# tzdata устанавливается в образ билдера
RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .

ADD go.sum .

# скачаиваются зависимости
RUN go mod download

# Копирование файлов проекта
COPY . .

#Билдим приложение (Удалено сообщение отладки -ldflags="-s -w" для уменьшения размера образа)
RUN go build -ldflags="-s -w" -o GoStats

## Deploy
FROM alpine

#Установлено ca-certificates, чтобы не было проблем с использованием TLS сертификатов.

RUN apk update --no-cache && apk add --no-cache ca-certificates

#tzdata устанавливается в builder-образ билдера, а в окончательный образ копируется только необходимый часовой пояс.
COPY --from=builder /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow

ENV TZ Europe/Moscow

WORKDIR /app

COPY --from=builder /build/GoStats /app/GoStats
COPY --from=builder /build/front /app/front

CMD ["./GoStats"]
