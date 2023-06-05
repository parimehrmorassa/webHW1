#get_users_with_sql_inject-service
FROM golang:1.19-alpine

WORKDIR /app
COPY . .

RUN go build -o get_users_with_sql_inject get_users_with_sql_inject.go

EXPOSE 50053

CMD ["./get_users_with_sql_inject"]
