package client

import (
	"context"
	"io/ioutil"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/RobertGrantEllis/t9/proto"
)

// Client is a simple interface for invoking T9 RPC methods.
type Client interface {
	proto.T9Client
	SimpleLookup(string, bool) ([]string, error)
}

// New instantiates a Client with the designated Configuration file, if any.
func New(configurations ...Configuration) (Client, error) {

	var configuration Configuration

	switch len(configurations) {
	case 0:
		configuration = NewConfiguration()
	case 1:
		configuration = configurations[0]
	default:
		return nil, errors.New(`cannot instantiate client with more than one configuration`)
	}

	if err := configuration.normalize(); err != nil {
		return nil, errors.Wrap(err, `cannot instantiate client`)
	}

	tlsConfig, err := getClientTLSConfig()
	if err != nil {
		return nil, err
	}

	if len(configuration.TrustedCertificateFile) > 0 {

		certificateBytes, err := ioutil.ReadFile(configuration.TrustedCertificateFile)
		if err != nil {
			return nil, errors.Wrap(err, `could not read trusted certificate file`)
		}

		tlsConfig.RootCAs.AppendCertsFromPEM(certificateBytes)
	}

	creds := credentials.NewTLS(tlsConfig)

	ctx, cancel := context.WithTimeout(context.Background(), configuration.ConnectionTimeout)
	defer cancel()

	options := append(configuration.DialOptions, grpc.WithTransportCredentials(creds))

	conn, err := grpc.DialContext(ctx, configuration.Address, options...)
	if err != nil {
		return nil, errors.Wrap(err, `could not connect to server`)
	}

	return &client{
		T9Client:      proto.NewT9Client(conn),
		Configuration: &configuration,
	}, nil
}

// NewDefaultClient returns a Client with default settings
func NewDefaultClient() Client {

	client, err := New()
	if err != nil {
		// should never happen
		panic(errors.Wrap(err, `cannot instantiate default client`))
	}

	return client
}

type client struct {
	proto.T9Client
	*Configuration
}

// SimpleLookup returns all word matches for the designated parameters.
func (c *client) SimpleLookup(digits string, exact bool) ([]string, error) {

	request := &proto.LookupRequest{
		Digits: digits,
		Exact:  exact,
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.RequestTimeout)
	defer cancel()

	response, err := c.T9Client.Lookup(ctx, request, grpc.FailFast(false))
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !response.Status {
		return nil, errors.New(response.Message)
	}

	return response.Words, nil
}
