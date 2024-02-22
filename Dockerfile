# Step 1: Build the application
# Use the official Golang 1.21 image as the base for the build stage
FROM golang:1.21 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Step 2: Use build output from 'builder' stage
# Start from a smaller image to keep the final image size down
FROM alpine:latest

# Add ca-certificates in case you need them
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy the encrypted.bin file into the container
COPY --from=builder /app/encrypted.bin .

# Set the environment variable to the path of encrypted.bin inside the container
ENV ENCRYPTED_CREDENTIALS_PATH="/root/encrypted.bin"

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
