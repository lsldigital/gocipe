package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"text/template"

	rice "github.com/GeertJohan/go.rice"
	"github.com/jinzhu/inflection"
)

// NoHeaderFormat is a special string to indicate no header in generated files
const NoHeaderFormat = "---"

var (
	templates *rice.Box

	// ErrorSkip is a pseudo-error to indicate code generation is skipped
	ErrorSkip = errors.New("skipped generation")

	reMatchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	reMatchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")

	templatesFn template.FuncMap

	// AppName represents the application name
	AppName string

	// AppImportPath represents GO import path for the project at hand
	AppImportPath string

	// WorkingDir represents current working directory
	WorkingDir string
)

func init() {
	templatesFn = template.FuncMap{
		"plus1": func(index int) int {
			return index + 1
		},
		"widget_field": func(component string, widget string, field Field) (string, error) {
			if widget == "" {
				return "", nil
			}
			return ExecuteTemplate(component+"_editor-field-"+widget+".vue.tmpl", field)
		},
		"plural": func(str string) string {
			return inflection.Plural(str)
		},
		"lower": func(str string) string {
			return strings.ToLower(str)
		},
		"upper": func(str string) string {
			return strings.ToUpper(str)
		},
		"ucfirst": func(str string) string {
			switch len(str) {
			case 0:
				return ""
			case 1:
				return strings.ToUpper(str)
			default:
				return strings.ToUpper(string(str[0])) + strings.ToLower(str[1:])
			}
		},
		"trimPrefix": func(str, prefix string) string {
			return strings.TrimPrefix(str, prefix)
		},
		"DerefAdminOpts": func(b *AdminOpts) (AdminOpts, error) {
			if b == nil {
				return AdminOpts{}, errors.New("admin opts is nil")
			}
			return *b, nil
		},
		"RelFuncName":          RelFuncName,
		"snake":                ToSnakeCase,
		"pkeyPropertyType":     GetPrimaryKeyDataType,
		"pkeyPropertyEmptyVal": GetPrimaryKeyEmptyVal,
		"pkeyIsAuto":           GetPrimaryKeyDataIsAuto,
		"pkeyIsInt":            GetPrimaryKeyDataIsInt,
		"pkeyFieldType":        GetPrimaryKeyFieldType,
		"fkeyPropertyTypeName": func(entities map[string]Entity, rel Relationship) string {
			return entities[rel.Entity].Name
		},
		"fkeyPropertyType": func(entities map[string]Entity, rel Relationship) (string, error) {
			return GetPrimaryKeyDataType(entities[rel.Entity].PrimaryKey)
		},
		// getAdminFilters returns possible filters from fields in an entity
		"getAdminFilters": func(entities map[string]Entity, entity Entity) AdminFilters {
			var (
				filters                                 AdminFilters
				filtersBool, filtersString, filtersDate []string
			)

			for _, field := range entity.Fields {
				switch field.Property.Type {
				case "bool":
					filtersBool = append(filtersBool, field.Schema.Field)
					filters.HasBool = true
				case "string":
					filtersString = append(filtersString, field.Schema.Field)
					filters.HasString = true
				case "date":
					filtersDate = append(filtersDate, field.Schema.Field)
					filters.HasDate = true
				}
			}

			for _, rel := range entity.Relationships {
				switch rel.Type {
				case RelationshipTypeOneOne, RelationshipTypeManyOne:
					switch entities[entity.Name].PrimaryKey {
					case PrimaryKeyUUID, PrimaryKeyString:
						filtersString = append(filtersString, rel.ThisID)
						filters.HasString = true
					}
				}
			}

			if len(filtersBool) != 0 {
				filters.BoolFilters = `"` + strings.Join(filtersBool, `","`) + `"`
			}

			if len(filtersString) != 0 {
				filters.StringFilters = `"` + strings.Join(filtersString, `","`) + `"`
			}

			if len(filtersDate) != 0 {
				filters.DateFilters = `"` + strings.Join(filtersDate, `","`) + `"`
			}

			return filters
		},
		"json": func(item interface{}) (string, error) {
			jsob, err := json.Marshal(item)
			if err == nil {
				return string(jsob), err
			}
			return "", err
		},
		"hasFileFields": func(entity Entity) bool {
			for _, field := range entity.Fields {
				switch field.EditWidget.Type {
				case WidgetTypeFile, WidgetTypeImage:
					return true
				}
			}
			return false
		},
		"getFileFields": func(entity Entity) []FileField {
			var fileFields []FileField
			for _, field := range entity.Fields {
				switch field.EditWidget.Type {
				case WidgetTypeFile, WidgetTypeImage:
					fileFields = append(fileFields, FileField{EntityName: entity.Name, PropertyName: field.Property.Name, FieldName: field.Schema.Field})
				}
			}

			return fileFields
		},
		"getEntityLabelField": func(entities map[string]Entity, name string) (string, error) {
			if entity, ok := entities[name]; ok {
				return entity.LabelField, nil
			}

			return "", errors.New("entity not found: " + name)
		},
	}
}

// GenerationWork represents generation work
type GenerationWork struct {
	Waitgroup *sync.WaitGroup
	Done      chan GeneratedCode
}

// GeneratedCode represents code that has been generated and the intended file
type GeneratedCode struct {
	// Generator indicates which generator produced the code
	Generator string

	// Filename is the name of the file to write to
	Filename string

	// Code is the generated code
	Code string

	// NoOverwrite will not overwrite an existing file
	NoOverwrite bool

	// Aggregate indicates file output needs to be aggregated instead of individual file write
	Aggregate bool

	// Error represents any error that may have occurred
	Error error

	// GeneratedHeaderFormat is used to prepend generated warning header on non-overwritable files. default: // %s
	GeneratedHeaderFormat string
}

// FileField used for getFileFields
type FileField struct {
	EntityName string
	PropertyName string
	FieldName string
}

// AdminFilters used for List
type AdminFilters struct {
	HasBool, HasString, HasDate             bool
	BoolFilters, StringFilters, DateFilters string
}

// SetTemplates injects template box
func SetTemplates(box *rice.Box) {
	templates = box
}

//GetAbsPath returns absolute path
func GetAbsPath(src string) (string, error) {
	gopath := os.Getenv("GOPATH")
	location := strings.TrimRight(strings.Replace(src, "$GOPATH", gopath, -1), "\\/")

	if !path.IsAbs(location) {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		location = path.Clean(wd + "/" + location)
	}

	return location, nil
}

// FileExists returns true is f (absolute path) exists; false otherwise
func FileExists(f string) bool {
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// DeleteIfExists deletes a file if exists
func DeleteIfExists(f string) {
	f, err := GetAbsPath(f)
	if err == nil && FileExists(f) {
		os.Remove(f)
	}
}

// ExecuteTemplate applies templating a text/template template given data and returns the string output
func ExecuteTemplate(name string, data interface{}) (string, error) {
	var output bytes.Buffer

	raw, err := templates.String(name)

	if err != nil {
		return "", err
	}

	tpl, err := template.New(name).Funcs(templatesFn).Parse(raw)
	if err != nil {
		return "", err
	}

	err = tpl.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

// ToSnakeCase converts a string to snake case
func ToSnakeCase(str string) string {
	snake := reMatchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = reMatchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// GetPrimaryKeyDataType returns golang data type for a given PrimaryKey type
func GetPrimaryKeyDataType(str string) (string, error) {
	switch str {
	case PrimaryKeySerial:
		return "int64", nil
	case PrimaryKeyInt:
		return "int64", nil
	case PrimaryKeyUUID:
		return "string", nil
	case PrimaryKeyString:
		return "string", nil
	}
	return "", errors.New("invalid primary key type: " + str)
}

// GetPrimaryKeyEmptyVal returns empty value for primary key
func GetPrimaryKeyEmptyVal(str string) (string, error) {
	switch str {
	case PrimaryKeySerial:
		return "0", nil
	case PrimaryKeyInt:
		return "0", nil
	case PrimaryKeyUUID:
		return `""`, nil
	case PrimaryKeyString:
		return `""`, nil
	}
	return "", errors.New("invalid primary key type: " + str)
}

// GetPrimaryKeyFieldType returns sql field data type for a given PrimaryKey type
func GetPrimaryKeyFieldType(str string) (string, error) {
	switch str {
	case PrimaryKeySerial:
		return "SERIAL", nil
	case PrimaryKeyInt:
		return "UNSIGNED INT", nil
	case PrimaryKeyUUID:
		return "CHAR(36)", nil
	case PrimaryKeyString:
		return "VARCHAR(255)", nil
	}
	return "", errors.New("invalid primary key type: " + str)
}

// GetPrimaryKeyDataIsAuto returns true if primary key is autogenerated
func GetPrimaryKeyDataIsAuto(str string) (bool, error) {
	switch str {
	case PrimaryKeySerial:
		return true, nil
	case PrimaryKeyInt:
		return false, nil
	case PrimaryKeyUUID:
		return true, nil
	case PrimaryKeyString:
		return false, nil
	}
	return false, errors.New("invalid primary key type: " + str)
}

// GetPrimaryKeyDataIsInt returns true if primary key is integer
func GetPrimaryKeyDataIsInt(str string) (bool, error) {
	switch str {
	case PrimaryKeySerial:
		return true, nil
	case PrimaryKeyInt:
		return true, nil
	case PrimaryKeyUUID:
		return false, nil
	case PrimaryKeyString:
		return false, nil
	}
	return false, errors.New("invalid primary key type: " + str)
}

// RelFuncName returns function name for a relationship repository function
func RelFuncName(rel Relationship) string {
	return inflection.Plural(strings.Title(strings.ToLower(rel.Name)))
}
