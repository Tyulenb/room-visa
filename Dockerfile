FROM golang:alpine

COPY . /app
WORKDIR /app

RUN go build -o cmd cmd/main.go
CMD ["/app/cmd/main"]
