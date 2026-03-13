FROM golang:1.22-alpine
RUN apk add --no-cache gcc musl-dev sqlite
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o music-api .
EXPOSE 24347
CMD ["./music-api"]