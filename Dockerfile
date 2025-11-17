FROM golang:1.23-alpine AS development

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["go", "run", "./cmd/result"]

FROM development AS builder

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /result ./cmd/result && \
    CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /frequentNumbers ./cmd/frequentNumbers && \
    CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /numberCount ./cmd/numberCount && \
    CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /consecutivePattern ./cmd/consecutivePattern

FROM alpine:latest AS production

RUN apk --no-cache add ca-certificates tzdata && \
    update-ca-certificates

WORKDIR /app

COPY --from=builder /result .
COPY --from=builder /frequentNumbers .
COPY --from=builder /numberCount .
COPY --from=builder /consecutivePattern .

RUN adduser -D -s /bin/sh appuser && \
    chown -R appuser:appuser /app
USER appuser

CMD ["./result"]