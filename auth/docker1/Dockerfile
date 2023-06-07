# image: req_pq-service

FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod go.sum ./ 

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

# RUN go build -o req_pq req_pq.go

EXPOSE 50052

CMD ["/docker-gs-ping"]


