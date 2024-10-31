FROM golang:alpine
WORKDIR /
COPY . .
RUN go build -o orders -mod vendor ./cmd/main.go
CMD ["./orders"]