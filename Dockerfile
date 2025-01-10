# Use the official Go image as the base image
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum into the container
COPY go.mod ./

# Download dependencies
RUN go mod tidy

# Copy the entire project into the container
COPY . .

# Build the Go application
RUN go build -o main ./cmd

# Use a lightweight image for the final container
FROM debian:bookworm-slim

# Set the working directory for the runtime container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8888

# Command to run the application
CMD ["./main"]
