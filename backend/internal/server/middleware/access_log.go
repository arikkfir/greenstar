package middleware

import (
	"bytes"
	"github.com/arikkfir/greenstar/backend/internal/server/util"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func AccessLogMiddleware(logSuccessfulRequests, excludeRemoteAddr bool, excludedHeaderPatterns []string, next http.Handler) http.Handler {
	excludedHeadersPatterns := make([]regexp.Regexp, len(excludedHeaderPatterns))
	for i, pattern := range excludedHeaderPatterns {
		excludedHeadersPatterns[i] = *regexp.MustCompile(pattern)
	}
	includeHeader := func(name string) bool {
		for _, pattern := range excludedHeadersPatterns {
			if pattern.MatchString(name) {
				return false
			}
		}
		return true
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := util.Logger(r.Context())

		// Add common request data
		accessLogLogger := l.With(
			"http:req:host", r.Host,
			"http:req:method", r.Method,
			"http:req:proto", r.Proto,
			"http:req:requestURI", r.RequestURI,
		)

		if !excludeRemoteAddr {
			accessLogLogger = accessLogLogger.With("http:req:remoteAddr", r.RemoteAddr)
		}

		// Add transfer encoding
		if len(r.TransferEncoding) > 0 {
			accessLogLogger = accessLogLogger.With("http:req:transferEncoding", r.TransferEncoding)
		}

		// Add headers (excluding some)
		for name, values := range r.Header {
			name = strings.ToLower(name)
			if includeHeader(name) {
				accessLogLogger = accessLogLogger.With("http:req:header:"+name, values)
			}
		}

		// Add trailer headers (excluding some)
		for name, values := range r.Trailer {
			name = strings.ToLower(name)
			if includeHeader(name) {
				accessLogLogger = accessLogLogger.With("http:req:trailer:"+name, values)
			}
		}

		// Keep a copy of the request body
		requestBody := bytes.Buffer{}
		responseBody := bytes.Buffer{}
		origReqBody := r.Body
		defer func() { r.Body = origReqBody }()
		r.Body = io.NopCloser(io.TeeReader(r.Body, &requestBody))
		rr := &responseBodyRecorder{
			ResponseWriter: w,
			w:              io.MultiWriter(&responseBody, w),
			status:         200,
		}

		// Invoke & time the next handler
		start := time.Now()
		next.ServeHTTP(rr, r)
		if logSuccessfulRequests || rr.status >= 300 {
			duration := time.Since(start)

			// Add request & response bodies
			if rr.status >= 400 {
				if requestBody.Len() > 1024*4 {
					requestBody.Truncate(1024 * 4)
					requestBody.WriteString("...")
				}
				if responseBody.Len() > 1024*4 {
					responseBody.Truncate(1024 * 4)
					responseBody.WriteString("...")
				}
				accessLogLogger = accessLogLogger.With("http:req:body", requestBody.Bytes())
				accessLogLogger = accessLogLogger.With("http:res:body", string(responseBody.Bytes()))
			}

			// Add invocation result
			accessLogLogger = accessLogLogger.With("http:process:duration", duration)
			accessLogLogger = accessLogLogger.With("http:res:status", rr.status)
			accessLogLogger = accessLogLogger.With("http:res:size", responseBody.Len())

			// Add response headers
			for name, values := range w.Header() {
				name := strings.ToLower(name)
				if includeHeader(name) {
					accessLogLogger = accessLogLogger.With("http:res:header:"+name, values)
				}
			}

			// Perform the logging with all the information we've added so far
			const message = "HTTP Request processed"
			if rr.status < 400 {
				accessLogLogger.InfoContext(r.Context(), message)
			} else if rr.status < 500 {
				accessLogLogger.WarnContext(r.Context(), message)
			} else {
				accessLogLogger.ErrorContext(r.Context(), message)
			}
		}
	})
}
