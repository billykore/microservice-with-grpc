FROM golang:1.19-alpine AS build

WORKDIR $GOPATH/src/microservice-with-grpc

COPY pb/payment/v1 ./pb/payment/v1
COPY gen/payment/v1 ./gen/payment/v1
COPY pb/transaction/v1 ./pb/transaction/v1
COPY gen/transaction/v1 ./gen/transaction/v1
COPY grpc/payment ./grpc/payment
COPY database ./database
COPY entity ./entity
COPY internal ./internal

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go build -o bin/app ./grpc/payment

# Final stage
FROM alpine:latest

WORKDIR /

COPY --from=build /go/src/microservice-with-grpc/bin/app ./

EXPOSE 50053

ENTRYPOINT ["./app"]
