# Build Stage
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Run Stage
FROM golang:1.25.6-alpine

WORKDIR /app

# Install ca-certificates, netcat, postgresql-client, curl, migrate and golangci-lint
RUN apk --no-cache add ca-certificates netcat-openbsd postgresql-client curl \
    && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz \
    && mv migrate /usr/local/bin/migrate \
    && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.8.0

# Copy binary from builder
COPY --from=builder /app/main /usr/local/bin/app-server

# Copy migrations and scripts
COPY db/migrations ./db/migrations
COPY db/scripts ./db/scripts
COPY scripts ./scripts
COPY . .

# Expose port
EXPOSE 8080

# Command to run
CMD ["app-server"]
