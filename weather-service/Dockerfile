FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o weather-api ./cmd/api/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/weather-api .

EXPOSE 8080

CMD ["./weather-api"]
