FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
RUN go build -o /tars ./cmd/

FROM alpine:latest

RUN apk add --no-cache tzdata sqlite

WORKDIR /app

COPY --from=builder /tars /app/tars
COPY --from=builder /app/config /app/config
COPY --from=builder /app/internal/application/email/templates /app/internal/application/email/templates

ENTRYPOINT ["/app/tars"]
