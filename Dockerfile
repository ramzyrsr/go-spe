# Step 1: Build the Go application
FROM golang:1.20-alpine AS build

# Install dependencies
RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go app, specifying the path to the main.go inside the cmd folder
RUN go build -o server ./cmd/main.go 

# Step 2: Run the Go application
FROM alpine:latest

# Install necessary dependencies
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary from the previous stage
COPY --from=build /app/server .

# Expose port for your API
EXPOSE 8080

# Run the Go binary
CMD ["./server"]
