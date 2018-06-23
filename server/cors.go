package server

import (
	"net/http"
)

// This allows browser requests from everywhere. That's totally fine.
func enableCors(h http.Handler) http.Handler {

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		origin := req.Header.Get(`Origin`)

		if origin != `` {

			rw.Header().Set("Access-Control-Allow-Origin", origin)

			if req.Method == `OPTIONS` && req.Header.Get(`Access-Control-Request-Method`) != `` {
				rw.Header().Set(`Access-Control-Allow-Headers`, `Content-Type,Accept`)
				rw.Header().Set(`Access-Control-Allow-Methods`, http.MethodGet)
				return
			}
		}

		h.ServeHTTP(rw, req)
	})
}
