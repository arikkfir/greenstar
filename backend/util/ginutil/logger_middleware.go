package ginutil

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
)

func SetLoggerMiddleware(c *gin.Context) {
	origCtx := c.Request.Context()

	e := log.Ctx(origCtx).With().
		Str("http:req:host", c.Request.Host).
		Str("http:req:method", c.Request.Method).
		Str("http:req:proto", c.Request.Proto).
		Str("http:req:remoteAddr", c.Request.RemoteAddr).
		Str("http:req:requestURI", c.Request.RequestURI)

	transferEncoding := zerolog.Arr()
	for _, encoding := range c.Request.TransferEncoding {
		transferEncoding = transferEncoding.Str(encoding)
	}
	e = e.Array("http:req:transferEncoding", transferEncoding)

	for name, values := range c.Request.Header {
		name = strings.ToLower(name)
		if strings.HasPrefix(name, "sec-") {
			continue
		}
		arr := zerolog.Arr()
		for _, value := range values {
			arr.Str(value)
		}
		e = e.Array("http:req:header:"+name, arr)
	}

	for name, values := range c.Request.Trailer {
		name = strings.ToLower(name)
		if strings.HasPrefix(name, "sec-") {
			continue
		}
		arr := zerolog.Arr()
		for _, value := range values {
			arr.Str(value)
		}
		e = e.Array("http:req:trailer:"+name, arr)
	}

	logger := e.Logger()

	newContextWithReqLogger := logger.WithContext(origCtx)
	c.Request = c.Request.WithContext(newContextWithReqLogger)
	c.Next()
	c.Request = c.Request.WithContext(origCtx)
}
