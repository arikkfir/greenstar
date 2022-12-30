package ginutil

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mssola/user_agent"
	"github.com/rs/zerolog/log"
	"strconv"
)

var (
	APIOfferedContentTypes = []string{
		binding.MIMEJSON,
		binding.MIMEXML,
		binding.MIMEYAML,
	}
	errorBrowserOfferedContentTypes = []string{
		"text/html",
		"text/plain",
		"application/json",
		"application/xml",
		"application/yaml",
	}
	errorDefaultOfferedContentTypes = []string{
		"application/json",
		"application/yaml",
		"application/xml",
		"text/plain",
		"text/html",
	}
)

func GetOfferedContentTypesForRequest(c *gin.Context, browserOfferedContentTypes, defaultOfferedContentTypes []string) []string {
	ua := user_agent.New(c.Request.UserAgent())
	if ua != nil && !ua.Bot() {
		if browserName, _ := ua.Browser(); browserName != "" {
			return browserOfferedContentTypes
		}
	}
	return defaultOfferedContentTypes
}

func RenderError(c *gin.Context, code int, err error) {
	log.Ctx(c.Request.Context()).Error().Err(err).Msgf("Request handling encountered an error")
	c.Negotiate(code, gin.Negotiate{
		Offered: GetOfferedContentTypesForRequest(c, errorBrowserOfferedContentTypes, errorDefaultOfferedContentTypes),
		Data: map[string]interface{}{
			"errors": []string{err.Error()},
		},
		HTMLName: strconv.Itoa(code) + ".html",
	})
}
