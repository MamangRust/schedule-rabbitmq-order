FROM golang:1.23-alpine

COPY . /app

WORKDIR /app

RUN go mod download

RUN go build -o main .

CMD ["./main"] 