package util

import (
	"bytes"
	"errors"
	"os"
	"path"
	"regexp"
	"strings"
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
		"snake":                ToSnakeCase,
		"pkeyPropertyEmptyVal": GetPrimaryKeyEmptyVal,
		"pkeyIsAuto":           GetPrimaryKeyDataIsAuto,
		"fkeyPropertyTypeName": func(entities map[string]Entity, rel Relationship) string {
			return entities[rel.Entity].Name
		},
	}
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
