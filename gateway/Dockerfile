FROM golang:1.19
WORKDIR /app

COPY . .

RUN go build -o gateway gateway.go 

# EXPOSE 8080
EXPOSE 6433
CMD ["./gateway"]

