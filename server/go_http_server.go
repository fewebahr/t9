package server

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/RobertGrantEllis/t9/logger"
)

type goHttpServer interface {
	Serve(net.Listener) error
	Shutdown(context.Context) error
}

func (s *server) instantiateGoHttpServer() {

	s.httpServer = &http.Server{
		Handler:  http.HandlerFunc(s.handleHttpRequest),
		ErrorLog: s.logger.GetLogger(logger.ErrorLevel),
	}
}

func (s *server) handleHttpRequest(rw http.ResponseWriter, req *http.Request) {

	isGrpcRequest := req.ProtoMajor == 2 && strings.Contains(req.Header.Get(`Content-Type`), `application/grpc`)
	isRestfulApiRequest := strings.HasPrefix(req.URL.Path, `/api/`)

	handler := s.frontendHandler

	switch {
	case isGrpcRequest:
		handler = s.grpcHandler
	case isRestfulApiRequest:
		handler = s.restfulApiHandler
	}

	handler.ServeHTTP(rw, req)
}
