gen_proto: clean
	protoc --go_out=. \
	--go-grpc_out=. \
	-I=$(PWD) pb/*/*/*.proto

clean:
	rm -rf gen