FROM golang:1.21

WORKDIR /app
COPY . .
RUN go build -o /main main.go
EXPOSE 17733
ENTRYPOINT ["/main", "--humidity-out", "/lwee/humidity.txt"]
