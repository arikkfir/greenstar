package ginutil

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initGin(devMode bool) {
	gin.DefaultWriter = log.Logger.Level(zerolog.TraceLevel)
	gin.DefaultErrorWriter = log.Logger.Level(zerolog.ErrorLevel)
	if devMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

func NewGin(devMode bool) *gin.Engine {
	initGin(devMode)
	router := gin.New()
	router.ContextWithFallback = true
	router.MaxMultipartMemory = 8 << 20
	router.Use(SetLoggerMiddleware)
	router.Use(AccessLogMiddleware)
	return router
}
