FROM golang:1.21

WORKDIR /app
COPY . .
RUN go build -o /main main.go
ENTRYPOINT ["/main", "--out", "/lwee/out.txt"]
