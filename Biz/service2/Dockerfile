# Dockerfile for get_users_with_sql_inject service (service2)
FROM golang:1.16-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the service source code
COPY service2 /app

# Build the service
RUN go build -o get_users_with_sql_inject .

# Expose the port that the service listens on
EXPOSE 50053

# Run the service
CMD ["./get_users_with_sql_inject"]
