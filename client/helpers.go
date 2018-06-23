package client

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/pkg/errors"

	"github.com/RobertGrantEllis/t9/bindata"
)

func getClientTlsConfig() (*tls.Config, error) {

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

	certPEM, err := bindata.Asset(`cert.pem`)
	if err != nil {
		return nil, errors.Wrap(err, `could not read embedded certificate`)
	}

	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, errors.Wrap(err, `could not get system certificate pool`)
	}

	pool.AppendCertsFromPEM(certPEM)

	return pool, nil
}
