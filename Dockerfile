FROM golang:1.21.6 AS builder
LABEL authors="vano"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o url_shortener

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/url_shortener .

EXPOSE 8000

CMD ["./url_shortener"]