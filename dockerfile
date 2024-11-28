# Base image
FROM golang:1.20

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./main"]
