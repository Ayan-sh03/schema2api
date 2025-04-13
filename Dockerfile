# Dockerfile for Schema2API

## Build Stage
FROM golang:1.22.1 AS builder
WORKDIR /app

# Copy go.mod and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o schema2api

## Final Stage
FROM alpine:latest
WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /app/schema2api .

# Copy the example schema files (optional)
COPY *.json ./

# Expose port 8081
EXPOSE 8081

# Set the binary as the entrypoint
ENTRYPOINT ["./schema2api"]