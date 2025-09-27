FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o weather-bff ./cmd/api/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/weather-bff .

EXPOSE 8080

CMD ["./weather-bff"]
