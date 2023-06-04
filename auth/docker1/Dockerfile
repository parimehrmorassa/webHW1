# image: req_pq
# Dockerfile for the first service (req_pq.go)
FROM golang:1.16-alpine

WORKDIR /app

COPY . .

RUN go build -o req_pq req_pq.go

EXPOSE 50052

CMD ["./req_pq"]


