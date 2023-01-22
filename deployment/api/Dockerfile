FROM golang:1.19-alpine AS build

WORKDIR $GOPATH/src/microservice-with-grpc

COPY api ./api
COPY gen ./gen

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go build -o bin/app ./api/main.go

# Final stage
FROM alpine:latest

WORKDIR /

COPY --from=build /go/src/microservice-with-grpc/bin/app ./

EXPOSE 8080

ENTRYPOINT ["./app"]