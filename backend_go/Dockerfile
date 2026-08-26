# Build stage — MenUnderFire Go API
FROM golang:1.26.1-alpine3.23 AS builder

WORKDIR /app

RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go binary with version embedded
ARG APP_VERSION=0.0.0
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-X main.Version=${APP_VERSION}" \
    -o server ./cmd/api

# Runtime stage
FROM alpine:3.23

RUN apk add --no-cache ca-certificates

RUN addgroup -S appgroup && adduser -S appuser -G appgroup && \
    mkdir -p /app/logs && \
    chown -R appuser:appgroup /app/logs

COPY --from=builder /app/server /app/server

ENV SERVER_PORT=7001
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=postgres
ENV DB_NAME=menunderfire
ENV DB_SSL_MODE=disable
ENV AUTH_DOMAIN=
ENV AUTH_AUDIENCE=none
ENV AUTH_ISSUER=
ENV AUTH_JWKS_URL=

EXPOSE 7001

WORKDIR /app

USER appuser

CMD ["/app/server"]
