package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/pkg/errors"

	"github.com/RobertGrantEllis/t9/assets"
)

func getClientTLSConfig() (*tls.Config, error) {

	certPool, err := getCertPool()
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		RootCAs:    certPool,
		ServerName: `t9`,
	}, nil
}

func getCertPool() (*x509.CertPool, error) {

	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, errors.Wrap(err, `could not get system certificate pool`)
	}

	fmt.Printf("adding certificate to trusted pool:\n%s\n", assets.Cert)

	pool.AppendCertsFromPEM(assets.Cert)

	return pool, nil
}
