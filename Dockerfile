# Stage 1: Modules caching (Common base)
FROM golang:1.24-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Stage 2: Development (with Air for hot reload)
FROM base AS dev
COPY . .
CMD ["go", "tool", "air", "-c", ".air.toml"]

# Stage 3: Builder (Compiles the binary for production)
FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Stage 4: Production (Minimal Alpine image)
FROM alpine:latest AS prod
WORKDIR /app
# Copy the built binary from the builder image
COPY --from=builder /app/server .
CMD ["./server"]