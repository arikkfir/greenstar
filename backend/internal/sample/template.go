package sample

import (
	"github.com/Masterminds/sprig/v3"
	"text/template"
)

func newTemplate(name, s string) (*template.Template, error) {
	return template.New(name).Funcs(sprig.FuncMap()).Funcs(template.FuncMap{
		"dateModify": dateModify,
		"rfc3339":    rfc3339,
	}).Parse(s)
}
