# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod ./
# COPY go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o q8-agent ./cmd/agent/main.go

# Final stage
FROM alpine:3.19

# Install docker and docker-compose
RUN apk add --no-cache docker-cli docker-cli-compose sudo

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/q8-agent .

# Create tenants root
RUN mkdir -p /opt/tenants

# Configure environment defaults
ENV Q8_AGENT_PORT=8080
ENV Q8_AGENT_ADMIN_TOKEN=change-me
ENV Q8_TENANTS_ROOT=/opt/tenants

EXPOSE 8080

# The agent needs to talk to the host's docker socket
# This will be mapped in docker-compose.yml or docker run
CMD ["./q8-agent"]
