package server

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"gitlab.com/fubahr/pipe"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/fewebahr/t9/proto"
)

func (s *server) instantiateAndRegisterRestfulHandler(pipeListener pipe.Listener, tlsConfig *tls.Config) error {
	creds := credentials.NewTLS(tlsConfig)
	conn, err := grpc.DialContext(
		context.Background(),                 // not a real connection so a timeout will never happen and background context is fine
		"",                                   // no address is required; it is not applicable for in-memory pipe
		grpc.WithTransportCredentials(creds), // GRPC still needs TLS even though this pipe is in-memory
		grpc.WithContextDialer(pipeListener.DialContextGRPC), // encapsulates logic for reaching the in-memory pipe listener
	)
	if err != nil {
		// should never happen
		return fmt.Errorf(`could not dial to in-memory pipe listener: %w`, err)
	}

	// default settings are totally fine since this is an in-memory connection
	client := proto.NewT9Client(conn)

	restHandler := runtime.NewServeMux()
	err = proto.RegisterT9HandlerClient(
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
