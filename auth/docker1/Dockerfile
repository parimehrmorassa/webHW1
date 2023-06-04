# image: req_pq
FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY req_pq.go .

RUN go build -o req_pq .

CMD ["./req_pq"]


