# Stage 1: Build the application
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o secure-image-service ./cmd/server

# Stage 2: Create the final, smaller image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/secure-image-service .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./secure-image-service"]

