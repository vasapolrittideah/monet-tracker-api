FROM golang:1.24-alpine

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY ./ ./

RUN go mod download

WORKDIR /app/swagger

ENTRYPOINT ["air", "-c", ".air.toml"]
