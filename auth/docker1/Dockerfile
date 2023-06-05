# image: req_pq-service

FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go build -o req_pq req_pq.go

EXPOSE 50052

CMD ["./req_pq"]


