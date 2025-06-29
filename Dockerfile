# Simple Go application with tests
FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go test -v ./internal/api

RUN go build -o main ./cmd/main.go

EXPOSE 8080

ENV HOST=0.0.0.0
ENV PORT=8080

CMD ["./main"]
