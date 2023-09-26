package middleware

import "net/http"

type responseStatusCodeRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseStatusCodeRecorder) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
	r.statusCode = status
}
