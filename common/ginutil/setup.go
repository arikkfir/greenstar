package ginutil

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
)

//go:embed templates
var templates embed.FS

func NewRouter() (*gin.Engine, error) {
	router := gin.Default()
	router.SetHTMLTemplate(template.Must(template.ParseFS(templates, "templates/*.html")))
	router.MaxMultipartMemory = 8 << 20
	return router, nil
}
