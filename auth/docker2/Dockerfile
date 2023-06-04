# Use a golang base image
FROM golang:1.16-alpine

# Set the working directory
WORKDIR /app

# Copy the source code into the container
COPY req_DH_params.go .

# Build the Go application
RUN go build -o req_DH_params

# Set the entry point command
CMD ["./req_DH_params"]
