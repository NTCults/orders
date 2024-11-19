FROM golang:alpine
WORKDIR /
COPY . .
RUN go build -o orders -mod vendor ./cmd/orders/main.go
CMD ["./orders"]