package middleware

import "net/http"

type responseBodyDiscarder struct {
	http.ResponseWriter
}

func (r *responseBodyDiscarder) Write(b []byte) (int, error) {
	return 0, nil
}

func (r *responseBodyDiscarder) WriteHeader(status int) {
	r.ResponseWriter.WriteHeader(status)
}
