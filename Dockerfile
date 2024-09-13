FROM golang:1.23-alpine AS build

WORKDIR /app

COPY . /app
COPY local.env /app

RUN go mod download

RUN go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /root

COPY --from=build /app .

EXPOSE 8080

CMD ["./main"]