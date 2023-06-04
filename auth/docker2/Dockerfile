# Dockerfile for the second service (req_DH_params.go)
FROM golang:1.16-alpine

WORKDIR /app

COPY . .

RUN go build -o req_DH_params req_DH_params.go

EXPOSE 50054

CMD ["./req_DH_params"]

CMD ["./req_DH_params"]
