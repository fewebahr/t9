package server

import (
	"fmt"
	"net/http"

	"github.com/fewebahr/t9/logger"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(statusCode int) {
	if lrw.statusCode != 0 {
		return // status code already written
	}
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *loggingResponseWriter) Write(buffer []byte) (int, error) {
	if lrw.statusCode == 0 {
		lrw.statusCode = http.StatusOK
	}
	return lrw.ResponseWriter.Write(buffer)
}

func enableLoggingMiddleware(inner http.Handler, l logger.Logger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		lrw := &loggingResponseWriter{ResponseWriter: rw}
		inner.ServeHTTP(lrw, req)

		message := fmt.Sprintf("%d %s", lrw.statusCode, path)
		switch {
		case lrw.statusCode >= 500: // server error
			l.Error(message)
		case lrw.statusCode >= 400: // client error
			l.Warn(message)
		default: // 1xx, 2xx, and 3xx are all fine
			l.Info(message)
		}

	})
}
