FROM golang:1.21.6-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o url_shortener

FROM alpine:latest

ENV MONGO_URI="mongodb://mongo:27017"

COPY --from=build /app/url_shortener .

EXPOSE 8000

CMD ["./url_shortener"]