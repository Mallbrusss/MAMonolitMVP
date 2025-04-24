FROM golang:1.24.0 AS builder

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

# RUN go test -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -v -o service ./cmd/app

FROM alpine:latest AS build-release-stage

WORKDIR /root/

RUN apk upgrade --no-cache && apk add libc6-compat curl

COPY --from=builder /app/service /service

CMD ["/service"]


