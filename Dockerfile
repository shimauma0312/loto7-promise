# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Build applications 
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /result ./cmd/result && \
    CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /frequentNumbers ./cmd/frequentNumbers

# Runtime stage - using alpine for better compatibility
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy built binaries from builder stage
COPY --from=builder /result .
COPY --from=builder /frequentNumbers .

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Default command
CMD ["./result"]