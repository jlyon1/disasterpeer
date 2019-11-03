FROM golang:latest

RUN mkdir /app

WORKDIR /app

COPY . .

RUN go get

RUN go build *.go
CMD ["./main"]