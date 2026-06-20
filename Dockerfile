# =====================================================
# Stage 1: Build
# =====================================================
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache gcc musl-dev ca-certificates tzdata

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -trimpath -ldflags="-s -w -X main.version=1.0.0 -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ) -X main.gitCommit=$(git rev-parse --short HEAD 2>/dev/null || echo unknown)" -o UserServer ./cmd/

# =====================================================
# Stage 2: Runtime
# =====================================================
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata curl \
    && addgroup -g 1001 appuser \
    && adduser -D -u 1001 -G appuser appuser

WORKDIR /app

COPY --from=builder /build/UserServer .
COPY --from=builder /build/pkg/config/config.yaml ./pkg/config/config.yaml

RUN mkdir -p /app/log && chown -R appuser:appuser /app

USER appuser

EXPOSE 8080
EXPOSE 443
EXPOSE 9090

HEALTHCHECK --interval=15s --timeout=3s --start-period=10s --retries=3 \
    CMD curl -sf http://localhost:8080/healthz || exit 1

ENTRYPOINT ["./UserServer"]
