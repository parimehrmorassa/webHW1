# Dockerfile for get_users service (service1)
FROM golang:1.16-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the service source code
COPY service1 /app

# Build the service
RUN go build -o get_users .

# Expose the port that the service listens on
EXPOSE 50051


#dsf
# Run the service
CMD ["./get_users"]
