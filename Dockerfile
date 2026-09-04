# ---- Build stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files and download dependencies (cache-friendly)
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY pkg/ ./pkg/

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/api

# ---- Runtime stage ----
FROM alpine:3.20

RUN apk --no-cache add ca-certificates && adduser -D appuser

WORKDIR /app
COPY --from=builder /server /app/server

USER appuser

EXPOSE 8080

CMD ["/app/server"]
