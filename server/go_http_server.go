package server

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/RobertGrantEllis/t9/logger"
)

type goHTTPServer interface {
	Serve(net.Listener) error
	Shutdown(context.Context) error
}

func (s *server) instantiateGoHTTPServer() {

	s.httpServer = &http.Server{
		Handler:  http.HandlerFunc(s.handleHTTPRequest),
		ErrorLog: s.logger.GetLogger(logger.ErrorLevel),
	}
}

func (s *server) handleHTTPRequest(rw http.ResponseWriter, req *http.Request) {

	isGrpcRequest := req.ProtoMajor == 2 && strings.Contains(req.Header.Get(`Content-Type`), `application/grpc`)
	isRestfulAPIRequest := strings.HasPrefix(req.URL.Path, `/api/`)

	handler := s.frontendHandler

	switch {
	case isGrpcRequest:
		handler = s.grpcHandler
	case isRestfulAPIRequest:
		handler = s.restfulAPIHandler
	}

	handler.ServeHTTP(rw, req)
}
