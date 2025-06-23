FROM golang:1.24.3-alpine AS base

RUN apk add --no-cache git curl

RUN go install github.com/air-verse/air@latest

FROM golang:1.24.3-alpine

COPY --from=base /go/bin/air /usr/bin/air

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENTRYPOINT ["air", "-c", ".air.toml"]