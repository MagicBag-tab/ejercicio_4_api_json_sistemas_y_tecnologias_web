FROM golang:1.22-alpine

WORKDIR /app

CMD ["go", "run", "main.go"]
