package server

import (
	"compress/gzip"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/RobertGrantEllis/t9/assets"
)

func (s *server) instantiateFrontendHandler() error {
	// initialize embedded filesystem
	filesystem, err := fs.Sub(assets.Frontend, "frontend")
	if err != nil {
		return fmt.Errorf("could not instantiate filesystem from embedded frontend: %w", err)
	}

	// turn into http Handler
	// handles default filename (index.html), content-type headers, etc.
	frontendHandler := http.FileServer(http.FS(filesystem))

	// support compression
	frontendHandler = enableCompressionMiddleware(frontendHandler, gzip.DefaultCompression)

	// support logging
	frontendHandler = enableLoggingMiddleware(frontendHandler, s.logger)

	s.frontendHandler = frontendHandler
	return nil
}
