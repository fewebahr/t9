package server

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/RobertGrantEllis/t9/logger"
)

type Server interface {
	Start()
	Stop()
}

func New(configuration Configuration) (Server, error) {

	if err := configuration.normalize(); err != nil {
		return nil, errors.Wrap(err, `cannot instantiate service`)
	}

	return &server{
		configuration: configuration,
		logger:        logger.New(configuration.getLogLevel()),
	}, nil
}

type server struct {
	configuration     Configuration
	logger            logger.Logger
	grpcHandler       http.Handler
	restfulApiHandler http.Handler
	frontendHandler   http.Handler
	httpServer        goHttpServer
}

func (s *server) Start() {

	tlsConfig, err := s.getTlsConfig()
	if err != nil {
		s.logger.Error(err)
		return
	}

	listener, err := s.getListener(tlsConfig)
	if err != nil {
		s.logger.Error(err)
		return
	}

	err = s.instantiateAndRegisterGrpcHandler()
	if err != nil {
		s.logger.Error(err)
		return
	}

	err = s.instantiateAndRegisterRestfulHandler()
	if err != nil {
		s.logger.Error(err)
		return
	}

	err = s.instantiateFrontendHandler()
	if err != nil {
		s.logger.Error(err)
		return
	}

	s.instantiateGoHttpServer()

	s.logger.Infof(`server listening on https://%s/`, s.configuration.Address)
	if err := s.httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
		s.logger.Error(errors.Wrap(err, `server terminated unexpectedly`))
	} else {
		s.logger.Info(`server stopped`)
	}
}

func (s *server) Stop() {

	s.httpServer.Shutdown(context.Background())
}
