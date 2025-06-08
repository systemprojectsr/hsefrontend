FROM golang:1.24.0-alpine

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go build -o app ./cmd/api

EXPOSE 8080

CMD ["./app"]