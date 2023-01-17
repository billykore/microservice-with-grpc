FROM golang:1.19-alpine

WORKDIR $GOPATH/src/microservice-with-grpc

COPY pb/customer/v1 ./pb/customer/v1
COPY gen/customer/v1 ./gen/customer/v1
COPY grpc/customer ./grpc/customer
COPY database ./database

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go build -o bin/app ./grpc/customer

EXPOSE 50051

CMD ["./bin/app"]
