package server

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"path/filepath"
	"strings"
)

func enableCompression(inner http.Handler, level int) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		compressor := newResponseCompressor(rw, req, level)
		inner.ServeHTTP(compressor, req)
		compressor.Close()
	})
}

type responseCompressor struct {
	io.Writer
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier
	level int
}

func newResponseCompressor(rw http.ResponseWriter, req *http.Request, level int) *responseCompressor {

	if level < gzip.DefaultCompression || level > gzip.BestCompression {
		level = gzip.DefaultCompression
	}

	compressor := &responseCompressor{ResponseWriter: rw, level: level}

	// the following steps guarantee that if the prior ResponseWriter implemented the
	// following interfaces, the new compressor will as well.

	if hijacker, ok := rw.(http.Hijacker); ok {
		compressor.Hijacker = hijacker
	}

	if flusher, ok := rw.(http.Flusher); ok {
		compressor.Flusher = flusher
	}

	if closeNotifier, ok := rw.(http.CloseNotifier); ok {
		compressor.CloseNotifier = closeNotifier
	}

	compressor.Writer = rw // by default

	if compressionShouldBeApplied(req) {
		// the request should have compression applied based on the content-type/extension/etc.
		// so... this function sets the writer based on what the client supports (if anything)
		compressor.setWriter(rw, req)
	}

	return compressor
}

func (rc *responseCompressor) WriteHeader(code int) {

	rc.ResponseWriter.Header().Del("Content-Length")
	rc.ResponseWriter.WriteHeader(code)
}

func (rc *responseCompressor) Write(b []byte) (int, error) {

	h := rc.ResponseWriter.Header()
	if h.Get("Content-Type") == "" {
		h.Set("Content-Type", http.DetectContentType(b))
	}
	h.Del("Content-Length")

	return rc.Writer.Write(b)
}

type flusher interface {
	Flush() error
}

func (rc *responseCompressor) Close() {

	// Flush compressed data if compressor supports it.
	if f, ok := rc.Writer.(flusher); ok {
		f.Flush()
	}

	// Flush HTTP response.
	if rc.Flusher != nil {
		rc.Flusher.Flush()
	}

	if closer, ok := rc.Writer.(io.Closer); ok {
		closer.Close()
	}
}

func (rc *responseCompressor) setWriter(rw http.ResponseWriter, req *http.Request) {

	acceptEncoding := strings.ToLower(req.Header.Get(`Accept-Encoding`))
	encodings := strings.Split(acceptEncoding, `,`)

	gzipSupported := false
	deflateSupported := false

	for _, encoding := range encodings {

		switch strings.TrimSpace(encoding) {
		case `gzip`:
			gzipSupported = true
		case `deflate`:
			deflateSupported = true
		}
	}

	switch {
	case gzipSupported:
		rw.Header().Set("Content-Encoding", "gzip")
		rw.Header().Add("Vary", "Accept-Encoding")

		// never returns an error as long as the level is kosher
		rc.Writer, _ = gzip.NewWriterLevel(rw, rc.level)
	case deflateSupported:
		rw.Header().Set("Content-Encoding", "deflate")
		rw.Header().Add("Vary", "Accept-Encoding")

		// never returns an error as long as the level is kosher
		rc.Writer, _ = flate.NewWriter(rw, rc.level)
	}
}

func compressionShouldBeApplied(req *http.Request) bool {

	_, filename := filepath.Split(req.RequestURI)
	if len(filename) == 0 {
		// there is a trailing slash so index html will be returned, and that is compressible
		return true
	}

	extension := strings.TrimPrefix(filepath.Ext(filename), `.`)

	// all the compressibles (for now)
	switch strings.ToLower(extension) {
	case `html`:
		fallthrough
	case `htm`:
		fallthrough
	case `js`:
		fallthrough
	case `css`:
		fallthrough
	case `json`:
		return true
	}

	// none of the above? Don't compress it. Probably an image or something
	return false
}
