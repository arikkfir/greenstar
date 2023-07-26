package web

import (
	"bytes"
	"github.com/arik-kfir/greenstar/backend/util"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strings"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	w      io.Writer
	status int
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	return r.w.Write(b)
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func AccessLogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create the logger event which we will start adding request & response data to
		event := log.Ctx(r.Context()).With()

		// Add common request data
		event = event.
			Str("http:req:id", r.Header.Get(RequestIDHeaderName)).
			Str("http:req:host", r.Host).
			Str("http:req:method", r.Method).
			Str("http:req:proto", r.Proto).
			Str("http:req:remoteAddr", r.RemoteAddr).
			Str("http:req:requestURI", r.RequestURI)

		// Add transfer encoding
		if len(r.TransferEncoding) > 0 {
			transferEncoding := zerolog.Arr()
			for _, encoding := range r.TransferEncoding {
				transferEncoding = transferEncoding.Str(encoding)
			}
			event = event.Array("http:req:transferEncoding", transferEncoding)
		}

		// Add headers (excluding some)
		for name, values := range r.Header {
			name = strings.ToLower(name)
			arr := zerolog.Arr()
			for _, value := range values {
				arr.Str(value)
			}
			event = event.Array("http:req:header:"+name, arr)
		}

		// Add trailer headers (excluding some)
		for name, values := range r.Trailer {
			name = strings.ToLower(name)
			arr := zerolog.Arr()
			for _, value := range values {
				arr.Str(value)
			}
			event = event.Array("http:req:trailer:"+name, arr)
		}

		// Keep a copy of the request body
		requestBody := bytes.Buffer{}
		responseBody := bytes.Buffer{}
		r.Body = &util.ReaderCloser{Reader: io.TeeReader(r.Body, &requestBody)}
		responseRecorder := &responseRecorder{
			ResponseWriter: w,
			w:              io.MultiWriter(&responseBody, w),
			status:         200,
		}

		// Invoke & time the next handler
		start := time.Now()
		next(responseRecorder, r.WithContext(event.Logger().WithContext(r.Context())))
		duration := time.Since(start)

		// Add request & response bodies
		event = event.Bytes("http:req:body", requestBody.Bytes())
		event = event.Bytes("http:res:body", responseBody.Bytes())

		// Add invocation result
		event = event.Dur("http:process:duration", duration)
		event = event.Int("http:res:status", responseRecorder.status)
		event = event.Int("http:res:size", responseBody.Len())

		// Add response headers
		for name, values := range w.Header() {
			arr := zerolog.Arr()
			for _, value := range values {
				arr.Str(value)
			}
			event = event.Array("http:res:header:"+strings.ToLower(name), arr)
		}

		// Perform the logging with all the information we've added so far
		const message = "HTTP Request processed"
		logger := &([]zerolog.Logger{event.Logger()}[0])
		if responseRecorder.status >= 200 && responseRecorder.status <= 399 {
			logger.Info().Msg(message)
		} else if responseRecorder.status >= 400 && responseRecorder.status <= 499 {
			logger.Warn().Msg(message)
		} else {
			logger.Error().Msg(message)
		}
	}
}
