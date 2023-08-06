FROM golang:1.20-alpine

RUN apk update && \
    apk upgrade && \
    apk add bash git && \
    rm -rf /var/cache/apk/*

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go install github.com/cosmtrek/air@latest

ENV CGO_ENABLED=0 \
    GOOS=linux 

CMD ["air", "-c", ".air.toml"]
