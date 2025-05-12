package api

import (
	"net/http"
	"time"

	"github.com/timhugh/digitalvenue/common"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{ResponseWriter: w}
}

func (rw *responseWriterWrapper) Status() int {
	return rw.status
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func LoggingMiddleware(log common.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				start := time.Now()
				wrapped := wrapResponseWriter(w)

				defer func() {
					if err := recover(); err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						var status int
						if wrapped.wroteHeader {
							status = wrapped.status
						} else {
							status = http.StatusInternalServerError
						}
						log.Error("[%d] %s %s - %s", status, r.Method, r.URL.Path, time.Since(start))
					}
				}()

				next.ServeHTTP(wrapped, r)
				log.Info("[%d] %s %s - %s", wrapped.status, r.Method, r.URL.Path, time.Since(start))
			},
		)
	}
}
