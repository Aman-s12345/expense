# ---- Build stage ----
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the main server
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server .

# Build the migration tool
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/migrate ./cmd/migrate

# ---- Runtime stage ----
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/migrate .
COPY --from=builder /app/db/migrations ./db/migrations

EXPOSE 3000

# Run migrations then start server
CMD ["sh", "-c", "./migrate && ./server"]