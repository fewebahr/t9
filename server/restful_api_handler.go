package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"github.com/RobertGrantEllis/t9/client"
	"github.com/RobertGrantEllis/t9/proto"
)

func (s *server) instantiateAndRegisterRestfulHandler() error {

	// includes embedded certificate and SNI which is sufficient
	// no matter what other certificates are supported by the server
	client := client.NewDefaultClient()

	restHandler := runtime.NewServeMux()
	err := proto.RegisterT9HandlerClient(
		context.Background(),
		restHandler,
		client,
	)
	if err != nil {
		return err
	}

	s.restfulAPIHandler = enableCors(restHandler)
	return nil
}
