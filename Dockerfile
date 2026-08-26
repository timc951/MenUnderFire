# Multi-stage Dockerfile for MenUnderFireApp
# Runs both frontend (nginx) and backend (Go) in the same container

# ============================================
# Stage 1: Build Go Backend
# ============================================
FROM golang:1.26.1-alpine3.23 AS backend-builder

WORKDIR /app

# Install git for Go modules that might need it
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY backend_go/go.mod backend_go/go.sum ./
RUN go mod download

# Copy source code
COPY backend_go/ ./

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/api

# ============================================
# Stage 2: Build React Frontend
# ============================================
FROM node:24-alpine AS frontend-builder

WORKDIR /app

# Copy package.json only - lockfile has platform-specific bindings
COPY frontend/package.json ./
# Fresh install to get correct Linux native bindings
RUN npm install

# Copy source code
COPY frontend/ ./

# Build arguments for frontend environment variables
ARG VITE_API_URL=/api
ARG VITE_AUTH_PUBLISHABLE_KEY
ARG VITE_DEBUG_BUILD=false

# Set environment variables for build
ENV VITE_API_URL=$VITE_API_URL
ENV VITE_AUTH_PUBLISHABLE_KEY=$VITE_AUTH_PUBLISHABLE_KEY
ENV VITE_DEBUG_BUILD=$VITE_DEBUG_BUILD

# Build the frontend
RUN npm run build

# ============================================
# Stage 3: Runtime Image
# ============================================
FROM alpine:3.23

# Install nginx and supervisor
RUN apk add --no-cache nginx supervisor

# Create non-root user and necessary directories
RUN addgroup -S appgroup && adduser -S appuser -G appgroup && \
    mkdir -p /var/log/supervisor /run/nginx /app/logs && \
    chown -R appuser:appgroup /var/log/supervisor /run/nginx /app/logs /var/lib/nginx /var/log/nginx

# Copy Go binary from builder
COPY --from=backend-builder /app/server /app/server

# Copy frontend build from builder
COPY --from=frontend-builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY docker/nginx.conf /etc/nginx/nginx.conf

# Copy supervisor configuration
COPY docker/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Environment variables for Go backend (with defaults)
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

# Expose ports
# 8080 - nginx (frontend + API proxy)
EXPOSE 8080

WORKDIR /app

USER appuser

# Start supervisor which manages both nginx and the Go backend
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
