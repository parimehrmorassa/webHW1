# get_users-service
FROM golang:1.19-alpine
WORKDIR /app
COPY  . .

RUN go build -o get_users get_users.go

EXPOSE 50051

CMD ["./get_users"]
