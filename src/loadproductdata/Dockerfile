FROM golang:1.16-alpine
WORKDIR /app
COPY main.go .
COPY go.mod .
COPY go.sum .
COPY products.json .
RUN go build -o main .
CMD ["./main"]
