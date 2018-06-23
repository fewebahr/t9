package server

import (
	"net"
	"strings"

	"github.com/pkg/errors"

	"github.com/RobertGrantEllis/t9/logger"
)

type Configuration struct {
	LogLevel string

	Address string

	DictionaryFile string

	CertificateFile string
	KeyFile         string

	CacheSize int
}

func NewConfiguration() Configuration {

	return Configuration{
		LogLevel:  logLevelDefault,
		Address:   listenerAddressDefault,
		CacheSize: cacheSizeDefault,
	}
}

func (configuration *Configuration) normalize() error {

	if len(configuration.LogLevel) == 0 {
		configuration.LogLevel = logLevelDefault
	} else if _, err := logger.ParseLevel(configuration.LogLevel); err != nil {
		return errors.Wrap(err, `invalid log level`)
	}

	if len(configuration.Address) == 0 {
		configuration.Address = listenerAddressDefault
	} else if _, _, err := net.SplitHostPort(configuration.Address); err != nil {
		return errors.Wrap(err, `invalid address`)
	}

	if configuration.CacheSize == 0 {
		configuration.CacheSize = cacheSizeDefault
	} else if configuration.CacheSize < 0 {
		return errors.Errorf(`invalid cache size (received %d)`, configuration.CacheSize)
	}

	configuration.CertificateFile = strings.TrimSpace(configuration.CertificateFile)
	configuration.KeyFile = strings.TrimSpace(configuration.KeyFile)

	if len(configuration.CertificateFile) > 0 && len(configuration.KeyFile) == 0 {
		return errors.New(`cannot designate certificate file but not key file`)
	} else if len(configuration.CertificateFile) == 0 && len(configuration.KeyFile) > 0 {
		return errors.New(`cannot designate key file but not certificate file`)
	}

	return nil
}

func (configuration *Configuration) getLogLevel() logger.Level {

	level, _ := logger.ParseLevel(configuration.LogLevel)
	return level
}
