package middleware

import (
	"io"
	"net/http"
)

type responseBodyRecorder struct {
	http.ResponseWriter
	w      io.Writer
	status int
}

func (r *responseBodyRecorder) Write(b []byte) (int, error) {
	return r.w.Write(b)
}

func (r *responseBodyRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
