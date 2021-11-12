package client

import (
	"net"
	"strings"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Configuration is a structure encapsulating all configurable elements for a Client.
type Configuration struct {
	Address string

	ConnectionTimeout time.Duration
	RequestTimeout    time.Duration
	DialOptions       []grpc.DialOption

	TrustedCertificateFile string
}

// NewConfiguration returns a default Configuration.
func NewConfiguration() Configuration {

	return Configuration{
		Address:           serverAddressDefault,
		ConnectionTimeout: connectionTimeoutDefault,
		RequestTimeout:    requestTimeoutDefault,
		DialOptions:       nil,
	}
}

func (configuration *Configuration) normalize() error {

	if len(configuration.Address) == 0 {
		configuration.Address = serverAddressDefault
	} else if _, _, err := net.SplitHostPort(configuration.Address); err != nil {
		return errors.Wrap(err, `invalid address`)
	}

	if configuration.ConnectionTimeout == 0 {
		configuration.ConnectionTimeout = connectionTimeoutDefault
	} else if configuration.ConnectionTimeout < 0 {
		return errors.Errorf(`invalid connection timeout (received: %d)`, configuration.ConnectionTimeout)
	}

	if configuration.RequestTimeout == 0 {
		configuration.RequestTimeout = requestTimeoutDefault
	} else if configuration.RequestTimeout < 0 {
		return errors.Errorf(`invalid request timeout (received: %d)`, configuration.RequestTimeout)
	}

	configuration.TrustedCertificateFile = strings.TrimSpace(configuration.TrustedCertificateFile)

	return nil
}
