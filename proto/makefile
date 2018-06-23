# Requires protobuf compiler (https://github.com/google/protobuf/blob/master/README.md#protocol-compiler-installation)
# Command for installation on a Mac (with Homebrew):
# brew install protobuf

all: t9.pb.go t9.pb.gw.go

t9.pb.go:
	protoc -I . -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins:. *.proto
t9.pb.gw.go:
	protoc -I . -I ${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. *.proto

clean:
	rm -f *.go