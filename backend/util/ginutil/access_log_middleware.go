package ginutil

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"strings"
	"time"
)

type customReadCloser struct {
	r io.Reader
	c io.Closer
}

func (rc *customReadCloser) Read(p []byte) (int, error) {
	return rc.r.Read(p)
}

func (rc *customReadCloser) Close() error {
	return rc.c.Close()
}

func AccessLogMiddleware(c *gin.Context) {
	if c.Request.RequestURI == "/healthz" {
		c.Next()
		return
	}

	//
	// Set up a request logger, which:
	// - adds simple metadata fields
	// - replaces request body reader with a Tee reader which also logs the request body to a side copy
	// - adds all transfer-encoding headers
	// - adds all headers
	// - adds all trailers
	//
	requestBody := bytes.Buffer{}
	c.Request.Body = &customReadCloser{r: io.TeeReader(c.Request.Body, &requestBody), c: c.Request.Body}

	start := time.Now()
	c.Next()
	duration := time.Since(start)

	re := log.Ctx(c).With()
	re = re.
		Dur("http:process:duration", duration).
		Int("http:res:status", c.Writer.Status()).
		Int("http:res:size", c.Writer.Size())
	for name, values := range c.Writer.Header() {
		if strings.HasPrefix(name, "sec-") {
			continue
		}
		arr := zerolog.Arr()
		for _, value := range values {
			arr.Str(value)
		}
		re = re.Array("http:res:header:"+strings.ToLower(name), arr)
	}

	var errorsArr []error
	for _, err := range c.Errors {
		errorsArr = append(errorsArr, err.Err)
	}
	re = re.Errs("http:res:errors", errorsArr)

	logger := re.Logger()

	status := c.Writer.Status()
	var event *zerolog.Event
	if status >= 200 && status <= 399 {
		event = logger.Info()
	} else if status >= 400 && status <= 499 {
		event = logger.Warn()
	} else {
		event = logger.Error()

	}
	event.Msg("Request processed")
}
