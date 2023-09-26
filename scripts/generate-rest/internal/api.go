package internal

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"slices"
)

var (
	allowedScopes = []ModelScope{ModelScopeTenant, ModelScopeGlobal}
)

func readAPI(apiFile string) (*API, error) {
	apiFileBytes, err := os.ReadFile(apiFile)
	if err != nil {
		return nil, fmt.Errorf("failed reading API file '%s': %w", apiFile, err)
	}

	decoder := yaml.NewDecoder(bytes.NewReader(apiFileBytes))
	decoder.KnownFields(true)

	api := API{}
	if err := decoder.Decode(&api); err != nil {
		return nil, fmt.Errorf("failed decoding API file '%s': %w", apiFile, err)
	}

	for modelName, model := range api.Models {
		model.Name = modelName
		for propertyName, property := range model.Properties {
			property.Name = propertyName
			if goType, typeScriptType, oasType, oasTypeFormat, err := inferTypesFromAPIType(property.Type); err != nil {
				return nil, fmt.Errorf("failed to infer types for property '%s.%s': %w", modelName, propertyName, err)
			} else {
				property.GoType = goType
				property.TypeScriptType = typeScriptType
				property.OpenAPIType = oasType
				property.OpenAPITypeFormat = oasTypeFormat
			}
		}

		if model.List != nil {
			for i, fp := range model.List.FilterProperties {
				if fp.Type == "" {
					if _, ok := model.Properties[fp.Name]; ok {
						model.List.FilterProperties[i].Type = model.Properties[fp.Name].Type
						model.List.FilterProperties[i].GoType = model.Properties[fp.Name].GoType
						model.List.FilterProperties[i].TypeScriptType = model.Properties[fp.Name].TypeScriptType
						model.List.FilterProperties[i].OpenAPIType = model.Properties[fp.Name].OpenAPIType
						model.List.FilterProperties[i].OpenAPITypeFormat = model.Properties[fp.Name].OpenAPITypeFormat
						model.List.FilterProperties[i].Optional = model.Properties[fp.Name].Optional
					} else {
						return nil, fmt.Errorf("failed to find corresponding model property '%s.%s' for list filter property '%s'", modelName, fp.Name, fp.Name)
					}
				} else if goType, typeScriptType, oasType, oasTypeFormat, err := inferTypesFromAPIType(fp.Type); err != nil {
					return nil, fmt.Errorf("failed to infer types for list filter property '%s' in model '%s': %w", fp.Name, modelName, err)
				} else {
					model.List.FilterProperties[i].GoType = goType
					model.List.FilterProperties[i].TypeScriptType = typeScriptType
					model.List.FilterProperties[i].OpenAPIType = oasType
					model.List.FilterProperties[i].OpenAPITypeFormat = oasTypeFormat
				}
			}
		}

		if model.Get != nil {
			for i, p := range model.Get.Parameters {
				if p.Type == "" {
					if _, ok := model.Properties[p.Name]; ok {
						model.Get.Parameters[i].Type = model.Properties[p.Name].Type
						model.Get.Parameters[i].GoType = model.Properties[p.Name].GoType
						model.Get.Parameters[i].TypeScriptType = model.Properties[p.Name].TypeScriptType
						model.Get.Parameters[i].OpenAPIType = model.Properties[p.Name].OpenAPIType
						model.Get.Parameters[i].OpenAPITypeFormat = model.Properties[p.Name].OpenAPITypeFormat
						model.Get.Parameters[i].Optional = model.Properties[p.Name].Optional
					} else {
						return nil, fmt.Errorf("failed to find corresponding model property '%s.%s' for list filter property '%s'", modelName, p.Name, p.Name)
					}
				} else if goType, typeScriptType, oasType, oasTypeFormat, err := inferTypesFromAPIType(p.Type); err != nil {
					return nil, fmt.Errorf("failed to infer types for list filter property '%s' in model '%s': %w", p.Name, modelName, err)
				} else {
					model.Get.Parameters[i].GoType = goType
					model.Get.Parameters[i].TypeScriptType = typeScriptType
					model.Get.Parameters[i].OpenAPIType = oasType
					model.Get.Parameters[i].OpenAPITypeFormat = oasTypeFormat
				}
			}
		}
	}

	return &api, nil
}

type API struct {
	Metadata Metadata          `yaml:"metadata"`
	Models   map[string]*Model `yaml:"models"`
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
	Name       string               `yaml:"-"`
	Scope      ModelScope           `yaml:"scope"`
	Path       string               `yaml:"path"`
	Properties map[string]*Property `yaml:"properties"`
	Create     *CreateOperation     `yaml:"create"`
	List       *ListOperation       `yaml:"list"`
	Get        *GetOperation        `yaml:"get"`
	Merge      *MergeOperation      `yaml:"merge"`
	Update     *UpdateOperation     `yaml:"update"`
	Delete     *DeleteOperation     `yaml:"delete"`
}

type Property struct {
	Name              string `yaml:"name"`
	Type              string `yaml:"type"`
	Sortable          bool   `yaml:"sortable"`
	GoType            string `yaml:"goType"`
	TypeScriptType    string `yaml:"typeScriptType"`
	OpenAPIType       string `yaml:"oasType"`
	OpenAPITypeFormat string `yaml:"oasTypeFormat"`
	Optional          bool   `yaml:"optional"`
	ReadOnly          bool   `yaml:"readOnly"`
}

type Permission struct {
	Scope                  ModelScope `yaml:"scope"`
	Permission             string     `yaml:"permission"`
	TenantPathVariableName string     `yaml:"tenantPathVariableName"`
}

func (l *Permission) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias Permission
	var temp alias
	if err := unmarshal(&temp); err != nil {
		return err
	} else if !slices.Contains(allowedScopes, temp.Scope) {
		return fmt.Errorf("invalid permission scope: %s", l.Scope)
	}

	if temp.Scope == ModelScopeTenant && temp.TenantPathVariableName == "" {
		temp.TenantPathVariableName = "tenantID"
	}
	l.Scope = temp.Scope
	l.Permission = temp.Permission
	l.TenantPathVariableName = temp.TenantPathVariableName
	return nil
}

type CreateOperation struct {
	Permissions     []Permission `yaml:"permissions"`
	Transactional   bool         `yaml:"transactional"`
	AllowExplicitID bool         `yaml:"allowExplicitID"`
}

type ListOperation struct {
	Permissions      []Permission `yaml:"permissions"`
	Transactional    bool         `yaml:"transactional"`
	FilterProperties []Property   `yaml:"filterProperties"`
}

func (l *ListOperation) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var temp struct {
		Permissions      []Permission `yaml:"permissions"`
		Transactional    bool         `yaml:"transactional"`
		FilterProperties []yaml.Node  `yaml:"filterProperties"`
	}
	if err := unmarshal(&temp); err != nil {
		return err
	}
	var filterProperties []Property
	for _, node := range temp.FilterProperties {
		var item Property
		switch node.Kind {
		case yaml.ScalarNode:
			var name string
			if err := node.Decode(&name); err != nil {
				return err
			}
			item = Property{Name: name}
		case yaml.MappingNode:
			if err := node.Decode(&item); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unexpected YAML node kind in the `filterProperties` list: %v", node.Kind)
		}
		filterProperties = append(filterProperties, item)
	}
	l.Permissions = temp.Permissions
	l.Transactional = temp.Transactional
	l.FilterProperties = filterProperties
	return nil
}

type GetOperation struct {
	Permissions   []Permission `yaml:"permissions"`
	Parameters    []Property   `yaml:"parameters"`
	Transactional bool         `yaml:"transactional"`
}

func (l *GetOperation) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var temp struct {
		Permissions   []Permission `yaml:"permissions"`
		Parameters    []yaml.Node  `yaml:"parameters"`
		Transactional bool         `yaml:"transactional"`
	}
	if err := unmarshal(&temp); err != nil {
		return err
	}
	var parameters []Property
	for _, node := range temp.Parameters {
		var item Property
		switch node.Kind {
		case yaml.ScalarNode:
			var name string
			if err := node.Decode(&name); err != nil {
				return err
			}
			item = Property{Name: name}
		case yaml.MappingNode:
			if err := node.Decode(&item); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unexpected YAML node kind in the `parameters` list: %v", node.Kind)
		}
		parameters = append(parameters, item)
	}
	l.Permissions = temp.Permissions
	l.Transactional = temp.Transactional
	l.Parameters = parameters
	return nil
}

type MergeOperation struct {
	Permissions        []Permission `yaml:"permissions"`
	Transactional      bool         `yaml:"transactional"`
	RequiredProperties []string     `yaml:"requiredProperties"`
}

type UpdateOperation struct {
	Permissions        []Permission `yaml:"permissions"`
	Transactional      bool         `yaml:"transactional"`
	RequiredProperties []string     `yaml:"requiredProperties"`
}

type DeleteOperation struct {
	Permissions   []Permission `yaml:"permissions"`
	Transactional bool         `yaml:"transactional"`
}
