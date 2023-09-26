package middleware

import "net/http"

func PreventCachingMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Expires", "0")
		w.Header().Set("Pragma", "no-cache")
		next.ServeHTTP(w, r)
	}
}
