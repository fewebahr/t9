package server

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/RobertGrantEllis/t9/logger"
)

// Configuration encapsulates all configurable parameters of a Server
type Configuration struct {
	LogLevel string

	Address string

	DictionaryFile string

	CertificateFile string
	KeyFile         string

	CacheSize int
}

const (
	logLevelDefault        = `info`
	listenerAddressDefault = `127.0.0.1:4239`
	cacheSizeDefault       = 16 * 1024
)

// NewConfiguration instantiates a Configuration with default settings.
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
		return fmt.Errorf("invalid log level: %w", err)
	}

	if len(configuration.Address) == 0 {
		configuration.Address = listenerAddressDefault
	} else if _, _, err := net.SplitHostPort(configuration.Address); err != nil {
		return fmt.Errorf(`invalid address: %w`, err)
	}

	if configuration.CacheSize == 0 {
		configuration.CacheSize = cacheSizeDefault
	} else if configuration.CacheSize < 0 {
		return fmt.Errorf(`invalid cache size (received %d)`, configuration.CacheSize)
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
