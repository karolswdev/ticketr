# Multi-stage build for ticktr Jira Story Creator
# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

# Install git for go mod dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy go mod file
COPY go.mod ./

# Download dependencies (go.sum will be created if not present)
RUN go mod download || true

# Copy source code
COPY . .

# Build the binary with optimization flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o jira-story-creator \
    cmd/jira-story-creator/main.go

# Stage 2: Create minimal runtime image
FROM alpine:3.18

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -g 1000 -S appuser && \
    adduser -u 1000 -S appuser -G appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/jira-story-creator .

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Set entrypoint
ENTRYPOINT ["./jira-story-creator"]

# Default command (can be overridden)
CMD ["--help"]