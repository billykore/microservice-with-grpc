FROM golang:1.19-alpine AS build

WORKDIR $GOPATH/src/microservice-with-grpc

COPY pb/notification/v1 ./pb/notification/v1
COPY gen/notification/v1 ./gen/notification/v1
COPY grpc/notification ./grpc/notification
COPY database ./database
COPY entity ./entity
COPY internal ./internal

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go build -o bin/app ./grpc/notification

# Final stage
FROM alpine:latest

WORKDIR /

COPY --from=build /go/src/microservice-with-grpc/bin/app ./

EXPOSE 50053

ENTRYPOINT ["./app"]