package util

import (
	"bytes"
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

var (
	templates *rice.Box

	// ErrorSkip is a pseudo-error to indicate code generation is skipped
	ErrorSkip = errors.New("skipped generation")

	reMatchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	reMatchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")

	templatesFn template.FuncMap
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
		"snake": ToSnakeCase,
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

// Toolset represents go tools used by the generators
type Toolset struct {
	GoImports string
	GoFmt     string
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
