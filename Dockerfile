FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . ./
RUN CGO_ENABLED=0 go build -o bin/api main.go

FROM alpine:3.19 AS final
WORKDIR /app
COPY --from=builder /app/bin/api ./
ENTRYPOINT ["./api", "server"]
EXPOSE 9000