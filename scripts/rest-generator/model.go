package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

// TODO: add duration property type

type PropertyType struct {
	Name           string
	GoType         string
	TypeScriptType string
	OpenAPIType    string
	OpenAPIFormat  string
}

var (
	PropertyTypeString = PropertyType{
		Name:           "String",
		GoType:         "string",
		TypeScriptType: "string",
		OpenAPIType:    "string",
		OpenAPIFormat:  "",
	}
	PropertyTypeBoolean = PropertyType{
		Name:           "boolean",
		GoType:         "bool",
		TypeScriptType: "boolean",
		OpenAPIType:    "boolean",
		OpenAPIFormat:  "",
	}
	PropertyTypeInteger = PropertyType{
		Name:           "integer",
		GoType:         "int64",
		TypeScriptType: "number",
		OpenAPIType:    "number",
		OpenAPIFormat:  "int64",
	}
	PropertyTypeDecimal = PropertyType{
		Name:           "decimal",
		GoType:         "decimal.Decimal",
		TypeScriptType: "number",
		OpenAPIType:    "number",
		OpenAPIFormat:  "double",
	}
	PropertyTypeTimestamp = PropertyType{
		Name:           "timestamp",
		GoType:         "time.Time",
		TypeScriptType: "Date",
		OpenAPIType:    "string",
		OpenAPIFormat:  "date-time",
	}
	PropertyTypeDate = PropertyType{
		Name:           "date",
		GoType:         "time.Time",
		TypeScriptType: "Date",
		OpenAPIType:    "string",
		OpenAPIFormat:  "date",
	}
	PropertyTypeTime = PropertyType{
		Name:           "time",
		GoType:         "time.Time",
		TypeScriptType: "Date",
		OpenAPIType:    "string",
		OpenAPIFormat:  "time",
	}

	TypesMap = map[string]PropertyType{
		"string":    PropertyTypeString,
		"boolean":   PropertyTypeBoolean,
		"integer":   PropertyTypeInteger,
		"decimal":   PropertyTypeDecimal,
		"timestamp": PropertyTypeTimestamp,
		"date":      PropertyTypeDate,
		"time":      PropertyTypeTime,
	}
)

type API struct {
	Metadata Metadata         `yaml:"metadata"`
	Models   map[string]Model `yaml:"models"`
}

func (api *API) UnmarshalYAML(value *yaml.Node) error {
	var raw struct {
		Metadata Metadata         `yaml:"metadata"`
		Models   map[string]Model `yaml:"models"`
	}

	if err := value.Decode(&raw); err != nil {
		return fmt.Errorf("failed decoding API: %w", err)
	}

	for name := range raw.Models {
		model := raw.Models[name]
		model.Name = name
		raw.Models[name] = model
	}

	api.Metadata = raw.Metadata
	api.Models = raw.Models
	return nil
}

type Metadata struct {
	Contact     Contact `yaml:"contact"`
	DisplayName string  `yaml:"displayName"`
	Description string  `yaml:"description"`
	Server      string  `yaml:"server"`
	Version     string  `yaml:"version"`
}

type Contact struct {
	Email string `yaml:"email"`
}

type Model struct {
	Name       string              `yaml:"-"`
	Path       string              `yaml:"path"`
	Properties map[string]Property `yaml:"properties"`
	Create     *CreateOperation    `yaml:"create"`
	List       *ListOperation      `yaml:"list"`
	Get        *GetOperation       `yaml:"get"`
	Update     *UpdateOperation    `yaml:"update"`
	Delete     *DeleteOperation    `yaml:"delete"`
}

func (m *Model) UnmarshalYAML(value *yaml.Node) error {
	var raw struct {
		Path       string              `yaml:"path"`
		Properties map[string]Property `yaml:"properties"`
		Create     *CreateOperation    `yaml:"create"`
		List       *ListOperation      `yaml:"list"`
		Get        *GetOperation       `yaml:"get"`
		Update     *UpdateOperation    `yaml:"update"`
		Delete     *DeleteOperation    `yaml:"delete"`
	}
	if err := value.Decode(&raw); err != nil {
		return fmt.Errorf("failed decoding model: %w", err)
	}

	for propertyName := range raw.Properties {
		p := raw.Properties[propertyName]
		p.Name = propertyName
		raw.Properties[propertyName] = p
	}

	m.Path = raw.Path
	m.Properties = raw.Properties
	m.Create = raw.Create
	m.List = raw.List
	m.Get = raw.Get
	m.Update = raw.Update
	m.Delete = raw.Delete

	if raw.List != nil {
		raw.List.applyPropertyTypes(m)
	}

	return nil
}

type Property struct {
	Name     string       `yaml:"-"`
	Type     PropertyType `yaml:"type"`
	Sortable bool         `yaml:"sortable"`
	Required bool         `yaml:"required"`
	ReadOnly bool         `yaml:"readOnly"`
}

func (p *Property) UnmarshalYAML(value *yaml.Node) error {
	var raw struct {
		Type     string `yaml:"type"`
		Sortable bool   `yaml:"sortable"`
		Required bool   `yaml:"required"`
		ReadOnly bool   `yaml:"readOnly"`
	}

	if err := value.Decode(&raw); err != nil {
		return fmt.Errorf("failed decoding property: %w", err)
	}

	p.Type = TypesMap[raw.Type]
	p.Sortable = raw.Sortable
	p.Required = raw.Required
	p.ReadOnly = raw.ReadOnly
	return nil
}

type CreateOperation struct {
	Permissions     []string `yaml:"permissions"`
	AllowExplicitID bool     `yaml:"allowExplicitID"`
}

type ListOperation struct {
	Permissions []string                  `yaml:"permissions"`
	Filters     map[string]ListFilter     `yaml:"filters"`
	Parameters  map[string]QueryParameter `yaml:"parameters"`
	// TODO: validate filter and parameter names do not overlap
}

func (l *ListOperation) UnmarshalYAML(value *yaml.Node) error {
	var raw struct {
		Permissions []string                  `yaml:"permissions"`
		Filters     map[string]ListFilter     `yaml:"filters"`
		Parameters  map[string]QueryParameter `yaml:"parameters"`
	}
	if err := value.Decode(&raw); err != nil {
		return fmt.Errorf("failed decoding list operation: %w", err)
	}

	for key := range l.Filters {
		if _, exists := l.Parameters[key]; exists {
			// Cannot have a filter that is also a parameter (or vice versa)
			// No need to traverse both maps - enough to traverse one to find that out
			return fmt.Errorf("overlapping filter and parameter name: %s", key)
		}
	}

	l.Permissions = raw.Permissions
	l.Filters = raw.Filters
	l.Parameters = raw.Parameters
	return nil
}

func (l *ListOperation) applyPropertyTypes(m *Model) {
	for name := range l.Filters {
		f := l.Filters[name]
		if f.Type == nil {
			if p, ok := m.Properties[name]; ok {
				f.Type = &p.Type
				l.Filters[name] = f
			}
		}
	}
}

type ListFilter struct {
	Type     *PropertyType `yaml:"type"`
	Required bool          `yaml:"required"`
}

func (lf *ListFilter) UnmarshalYAML(value *yaml.Node) error {
	var raw struct {
		Type     *string `yaml:"type"`
		Required bool    `yaml:"required"`
	}
	if err := value.Decode(&raw); err != nil {
		return fmt.Errorf("failed decoding list operation filter: %w", err)
	} else if raw.Type != nil {
		if t, ok := TypesMap[*raw.Type]; ok {
			lf.Type = &t
		} else {
			return fmt.Errorf("invalid list filter type: %s", *raw.Type)
		}
	}
	lf.Required = raw.Required
	return nil
}

type GetOperation struct {
	Permissions []string                  `yaml:"permissions"`
	Parameters  map[string]QueryParameter `yaml:"parameters"`
}

type QueryParameter struct {
	Type     PropertyType `yaml:"type"`
	Required bool         `yaml:"required"`
}

func (q *QueryParameter) UnmarshalYAML(value *yaml.Node) error {
	var raw struct {
		Type     string `yaml:"type"`
		Required *bool  `yaml:"required"`
	}
	if err := value.Decode(&raw); err != nil {
		return fmt.Errorf("failed decoding query parameter: %w", err)
	} else if t, ok := TypesMap[raw.Type]; ok {
		q.Type = t
	} else {
		return fmt.Errorf("invalid query parameter type: %s", raw.Type)
	}

	if raw.Required == nil {
		q.Required = false
	} else {
		q.Required = *raw.Required
	}

	return nil
}

type UpdateOperation struct {
	Permissions []string `yaml:"permissions"`
}

type DeleteOperation struct {
	Permissions []string `yaml:"permissions"`
}
