//go:build tools
// +build tools

package tools

import (
	// protoc-gen-go generates go files from protobuf definitions
	_ "github.com/golang/protobuf/protoc-gen-go"
	// protoc-gen-grpc-gateway generates GRPC Gateway code useful for transparently proxying RESTful HTTP calls to GRPC service calls
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	// protoc-gen-swagger automatically generates swagger documentation from protobuf definitions
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
)
