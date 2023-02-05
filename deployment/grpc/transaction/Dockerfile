FROM golang:1.19-alpine AS build

WORKDIR $GOPATH/src/microservice-with-grpc

COPY pb/transaction/v1 ./pb/transaction/v1
COPY gen/transaction/v1 ./gen/transaction/v1
COPY grpc/transaction ./grpc/transaction
COPY database ./database
COPY entity ./entity
COPY internal ./internal

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go build -o bin/app ./grpc/transaction

# Final stage
FROM alpine:latest

WORKDIR /

COPY --from=build /go/src/microservice-with-grpc/bin/app ./

EXPOSE 50055

ENTRYPOINT ["./app"]