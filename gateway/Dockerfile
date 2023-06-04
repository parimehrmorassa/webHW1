# Use the official Golang image as the base image
FROM golang:1.16-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the remaining source code
COPY . .

# Build the Go application
RUN go build -o gateway .

# Expose the port on which the service will run
EXPOSE 8080

# Run the Go application
CMD ["./gateway"]
