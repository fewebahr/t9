package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"gitlab.com/fubahr/pipe"

	"github.com/RobertGrantEllis/t9/logger"
)

// Server is an interface with simple methods to start or stop a T9 server.
type Server interface {
	Start()
	Stop()
}

// New instantiates a server using the designated Configuration, or an error if the configuration is invalid.
func New(configuration Configuration) (Server, error) {
	if err := configuration.normalize(); err != nil {
		return nil, fmt.Errorf(`cannot instantiate service: %w`, err)
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
	restfulAPIHandler http.Handler
	frontendHandler   http.Handler
	httpServer        goHTTPServer
}

func (s *server) Start() {
	if err := s.composeAndStart(); err != nil {
		s.logger.Error(err)
	}
}

func (s *server) composeAndStart() error {
	// configure a TLS configuration, suitable for both server-side and client-side
	tlsConfig, err := s.getTLSConfig()
	if err != nil {
		return err
	}

	/*
		Instantiate all required listeners
	*/

	pipeListener := pipe.Listen() // in-memory pipes suitable for GRPC gateway

	// TCP listener for requests from frontend or other API consumers
	tcpListener, err := net.Listen(`tcp`, s.configuration.Address)
	if err != nil {
		return fmt.Errorf(`could not instantiate TCP listener: %w`, err)
	}

	// GRPC service must listen on both in-memory pipe and TCP listener (for RPC consumers)
	teeListener := pipe.TeeListener(pipeListener, tcpListener)

	// TLS is required for both
	tlsListener := tls.NewListener(teeListener, tlsConfig)
	defer tlsListener.Close() // will close all underlying listeners

	/*
		Instantiate all handlers for frontend, GRPC methods, and restful methods
	*/
	if err := s.instantiateFrontendHandler(); err != nil {
		return err
	}

	if err := s.instantiateAndRegisterGrpcHandler(); err != nil {
		return err
	}

	if err := s.instantiateAndRegisterRestfulHandler(pipeListener, tlsConfig); err != nil {
		return err
	}

	// Instantiate the HTTP server which will handle all types of requests
	s.instantiateGoHTTPServer()

	// Start serving the listeners
	s.logger.Infof(`server listening on https://%s/`, s.configuration.Address)
	if err := s.httpServer.Serve(tlsListener); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf(`server terminated unexpectedly: %w`, err)
	}

	s.logger.Info(`server stopped`)
	return nil
}

const serverShutdownTimeout = 10 * time.Second

func (s *server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Errorf("server shutdown failed after %s: force closing", serverShutdownTimeout)
		s.httpServer.Close()
	}
}
