FROM golang:1.19
WORKDIR /app

COPY . .

RUN go build -o gateway 

# EXPOSE 8080
EXPOSE 6433
CMD ["./gateway"]

