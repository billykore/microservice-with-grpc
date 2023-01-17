FROM golang:1.19-alpine

WORKDIR $GOPATH/src/microservice-with-grpc

COPY api ./api
COPY gen/customer/v1 ./gen/customer/v1

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go build -o bin/app ./api/main.go

EXPOSE 8080

CMD ["./bin/app"]