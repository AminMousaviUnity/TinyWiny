# Official Go image as a parent image
FROM golang:1.22 AS builder

# Set the working dir inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./

# Download all dependencies
RUN go mod tidy

# Copy the entire application code into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Use a lightweight image for the final container
FROM debian:bookworm-slim

# Set the working dir
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the port your app runs on
EXPOSE 8888

# Command to run the application
CMD ["./main"]
