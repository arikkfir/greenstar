package internal

import (
	"embed"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
	"regexp"
	"text/template"
)

var (
	//go:embed templates/*.tmpl
	templatesFS          embed.FS
	templates            *template.Template
	resourceServerTmplRE *regexp.Regexp
	resourceClientTmplRE *regexp.Regexp
)

func init() {
	strcase.ConfigureAcronym("ParentID", "parentID")
	pluralizeClient := pluralize.NewClient()
	functions := template.FuncMap{
		"fail":             func(format string, args ...any) (string, error) { return "", fmt.Errorf(format, args...) },
		"modelScopeGlobal": func() ModelScope { return ModelScopeGlobal },
		"modelScopeTenant": func() ModelScope { return ModelScopeTenant },
		"toCamelCase":      strcase.ToCamel,
		"toLowerCamelCase": strcase.ToLowerCamel,
		"toSnake":          strcase.ToSnake,
		"toScreamingSnake": strcase.ToScreamingSnake,
		"toKebab":          strcase.ToKebab,
		"toScreamingKebab": strcase.ToScreamingKebab,
		"isPlural":         pluralizeClient.IsPlural,
		"toPlural":         pluralizeClient.Plural,
		"isSingular":       pluralizeClient.IsSingular,
		"toSingular":       pluralizeClient.Singular,
		"pluralize":        pluralizeClient.Pluralize,
	}

	templates = template.New("rest").Funcs(sprig.FuncMap()).Funcs(functions)
	if _, err := templates.ParseFS(templatesFS, "templates/*.tmpl"); err != nil {
		panic(fmt.Errorf("failed to parse Go templates: %w", err))
	}
	resourceServerTmplRE = regexp.MustCompile(`^resource.([^.]+\.go)\.tmpl$`)
	resourceClientTmplRE = regexp.MustCompile(`^resource\.client\.ts\.tmpl$`)
}

func isResourceServerTemplate(templateName string) (string, bool) {
	matches := resourceServerTmplRE.FindStringSubmatch(templateName)
	if len(matches) == 2 {
		return matches[1], true
	} else {
		return "", false
	}
}

func isResourceClientTemplate(templateName string) bool {
	return resourceClientTmplRE.MatchString(templateName)
}
