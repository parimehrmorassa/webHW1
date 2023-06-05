# req_pq_dh_params-service
FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go build -o req_DH_params  

EXPOSE 50054

CMD ["./req_DH_params"]
