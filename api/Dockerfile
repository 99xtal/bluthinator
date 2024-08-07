# Use the official Golang image as the base image
FROM golang:1.20-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY src/ .

# Build the Go app
RUN go build -o main ./cmd/api/main.go

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./main"]