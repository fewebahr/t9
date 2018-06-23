package server

import (
	"crypto/tls"
	"net"

	"github.com/pkg/errors"
)

func (s *server) getListener(tlsConfig *tls.Config) (net.Listener, error) {

	listener, err := net.Listen(`tcp`, s.configuration.Address)
	if err != nil {
		return nil, errors.Wrap(err, `could not instantiate listener`)
	}

	listener = tls.NewListener(listener, tlsConfig)
	return listener, nil
}
