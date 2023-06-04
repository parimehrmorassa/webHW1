# image: req_pq
FROM golang:1.16-alpine

WORKDIR /app

COPY req_pq.go .

RUN go build -o req_pq .

CMD ["./req_pq"]


