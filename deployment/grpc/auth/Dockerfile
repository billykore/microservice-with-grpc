FROM golang:1.19-alpine

WORKDIR $GOPATH/src/microservice-with-grpc

COPY pb/auth/v1 ./pb/auth/v1
COPY gen/auth/v1 ./gen/auth/v1
COPY grpc/auth ./grpc/auth
COPY database ./database

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go build -o bin/app ./grpc/auth

EXPOSE 50052

CMD ["./bin/app"]