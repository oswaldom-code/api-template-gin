# syntax=docker/dockerfile:1

FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download && go mod verify
COPY . ./
RUN CGO_ENABLED=0 go build -o bin/api-template main.go

ENTRYPOINT ["/app/bin/api-template", "server"]
EXPOSE 9000
