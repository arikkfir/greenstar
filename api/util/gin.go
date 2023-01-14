package util

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

func SetLoggerMiddleware(c *gin.Context) {
	origCtx := c.Request.Context()

	requestLogger := log.Ctx(origCtx).With().
		Str("http:req:proto", c.Request.Proto).
		Str("http:req:method", c.Request.Method).
		Str("http:req:requestURI", c.Request.RequestURI).
		Str("http:req:host", c.Request.Host).
		Logger()
	// TODO: check if URL query is logged (should be part of RequestURI)
	for name, values := range c.Request.Header {
		arr := zerolog.Arr()
		for _, value := range values {
			arr.Str(value)
		}
		requestLogger = requestLogger.With().
			Array("http:req:header:"+strings.ToLower(name), arr).
			Logger()
	}

	ctxWithLogger := requestLogger.WithContext(origCtx)
	c.Request = c.Request.WithContext(ctxWithLogger)
	c.Next()
	c.Request = c.Request.WithContext(origCtx)
}

func AccessLogMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	duration := time.Since(start)

	zctx := log.Ctx(c.Request.Context()).With().
		Dur("duration", duration).
		Int("http:res:status", c.Writer.Status()).
		Int("http:res:size", c.Writer.Size())
	for name, values := range c.Writer.Header() {
		arr := zerolog.Arr()
		for _, value := range values {
			arr.Str(value)
		}
		zctx = zctx.Array("http:res:header:"+strings.ToLower(name), arr)
	}

	logger := zctx.Logger()
	logger.Info().Msg("Request processed")
}
