# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o knucklebones .

# Final stage
FROM alpine:latest

# Add ca certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/knucklebones .

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./knucklebones"]
