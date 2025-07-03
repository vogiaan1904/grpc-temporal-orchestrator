# Build stage
FROM golang:1.23.6-alpine AS builder

# Install ca-certificates (no need for git anymore)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and pre-generated proto files
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o main ./cmd/server

# Production stage
FROM alpine:latest AS production

# Install ca-certificates for SSL/TLS connections
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

WORKDIR /app

# Copy binary and config from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

# Copy proto files (if needed at runtime)
COPY --from=builder /app/protogen ./protogen

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app
USER appuser

# Order orchestrator typically uses a different port
EXPOSE 50055

CMD ["./main"]
