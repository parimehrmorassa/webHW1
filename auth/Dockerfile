# Dockerfile for the overall container

# Build the first service
FROM golang:1.16-alpine AS build-docker1
WORKDIR /app
COPY docker1/ .
RUN go build -o req_pq req_pq.go

# Build the second service
FROM golang:1.16-alpine AS build-docker2
WORKDIR /app
COPY docker2/ .
RUN go build -o req_DH_params req_DH_params.go

# Final container
FROM alpine:latest
WORKDIR /app
COPY --from=build-docker1 /app/req_pq .
COPY --from=build-docker2 /app/req_DH_params .

EXPOSE 50052
EXPOSE 50054

CMD ["./req_pq"]
