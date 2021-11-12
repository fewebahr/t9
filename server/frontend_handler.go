package server

import (
	"compress/gzip"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/RobertGrantEllis/t9/assets"
)

func (s *server) instantiateFrontendHandler() error {

	var frontendHandler http.Handler

	// initialize with core handler
	frontendHandler = http.HandlerFunc(s.handleFrontendRequest)

	// support compression
	frontendHandler = enableCompression(frontendHandler, gzip.DefaultCompression)

	s.frontendHandler = frontendHandler
	return nil
}

func (s *server) handleFrontendRequest(rw http.ResponseWriter, req *http.Request) {

	path := req.URL.Path

	frontendPath, err := getFrontendPath(path)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		s.logger.Warnf(`404 %s %s`, path, err)
		return
	}

	reader, err := getEmbeddedFrontendFileReader(frontendPath)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		s.logger.Warnf(`404 %s %s`, path, err)
		return
	}

	addContentTypeHeaderForPath(rw, frontendPath)
	io.Copy(rw, reader)
	s.logger.Infof(`200 %s`, path)
}

func getFrontendPath(path string) (string, error) {

	if len(path) == 0 || path == `/` {
		// default path
		path = `/index.html`
	}

	frontendPath := filepath.Clean(filepath.Join(`frontend`, path))
	if !strings.HasPrefix(frontendPath, `frontend/`) {
		// this was deliberately trying to escape the webroot jail!
		return ``, errors.New(`requested path was outside the webroot`)
	}

	return frontendPath, nil
}

func getEmbeddedFrontendFileReader(path string) (io.Reader, error) {

	data, err := assets.Frontend.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return data, nil
}

func addContentTypeHeaderForPath(rw http.ResponseWriter, path string) {

	extension := filepath.Ext(path)
	contentType := mime.TypeByExtension(extension)

	if len(contentType) > 0 {
		rw.Header().Set(`Content-Type`, contentType)
	}
}
