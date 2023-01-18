FROM golang:1.19-alpine AS build

WORKDIR $GOPATH/src/microservice-with-grpc

COPY pb/customer/v1 ./pb/customer/v1
COPY gen/customer/v1 ./gen/customer/v1
COPY grpc/customer ./grpc/customer
COPY database ./database

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go build -o bin/app ./grpc/customer

# Final stage
FROM alpine:latest

WORKDIR /

COPY --from=build /go/src/microservice-with-grpc/bin/app ./

EXPOSE 50051

ENTRYPOINT ["./app"]
