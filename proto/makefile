# Requires protobuf compiler (https://github.com/google/protobuf/blob/master/README.md#protocol-compiler-installation)
# Command for installation on a Mac (with Homebrew):
# brew install protobuf
# go install github.com/golang/protobuf/protoc-gen-go
# go install github.com/golang/protobuf/protoc-gen-grpc-gateway
# go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
all: t9.pb.go t9.pb.gw.go

t9.pb.go:
	protoc -I . --go_out=plugins=grpc:. *.proto
t9.pb.gw.go:
	protoc -I . --grpc-gateway_out=logtostderr=true:. *.proto

clean:
	rm -f *.go
