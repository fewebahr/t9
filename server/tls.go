package server

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/RobertGrantEllis/t9/assets"
	"github.com/pkg/errors"
)

func (s *server) getTLSConfig() (*tls.Config, error) {

	certPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, errors.Wrap(err, `could not get system certificate pool`)
	}

	certPool.AppendCertsFromPEM(assets.Cert)

	embeddedCertificate, err := getEmbeddedCertificate()
	if err != nil {
		return nil, err
	}

	// for now, assume no certificates are configured
	getCertificate := func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		// since this version of the function assumes no other certificates are
		// designated, simply return the embedded certificate
		return embeddedCertificate, nil
	}

	if len(s.configuration.CertificateFile) > 0 {
		// keyfile is also designated since the configuration normalization checks
		designatedCertificate, err := getDesignatedCertificate(
			s.configuration.CertificateFile,
			s.configuration.KeyFile,
		)
		if err != nil {
			return nil, err
		}

		// a certificate was designated, so override the getCertificate function
		getCertificate = func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {

			if hello.ServerName == `t9` {
				// this is our internal client SNI designation so use embedded cert
				return embeddedCertificate, nil
			}

			// otherwise use the designated cert
			return designatedCertificate, nil
		}
	}

	tlsConfig := &tls.Config{
		ServerName:               "t9",
		RootCAs:                  certPool,
		NextProtos:               []string{`h2`, `http/1.1`},
		GetCertificate:           getCertificate,
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	return tlsConfig, nil
}

func getDesignatedCertificate(certFile, keyFile string) (*tls.Certificate, error) {

	certBuf, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	keyBuf, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tlsCert, err := tls.X509KeyPair(certBuf, keyBuf)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &tlsCert, nil
}

func getEmbeddedCertificate() (*tls.Certificate, error) {

	tlsCert, err := tls.X509KeyPair(assets.Cert, assets.Key)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &tlsCert, nil
}
