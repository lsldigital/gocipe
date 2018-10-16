package output

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/fluxynet/gocipe/util"
	log "github.com/sirupsen/logrus"
)

var (
	_recipePath                 string
	_log, _gofiles              []string
	_tools                      toolset
	_written, _skipped, _failed int
	_verbose                    bool
)

const (
	logSuccess = "✅ [Ok]"
	logSkipped = "👻 [Skipped]"
	logError   = "❗️ [Error]"
	logInfo    = "🦋 [Info]"
)

// Output is the implementation
type Output struct {
	*log.Logger
	gofiles                   []string
	success, failure, skipped int
}

// toolset represents go tools used by the generators
type toolset struct {
	GoImports string
	GoFmt     string
	Protoc    string
	Dep       string
}

func init() {
	initToolset()
}

// SetVerbose can be used to toggle verbosity
func SetVerbose(verbose bool) {
	_verbose = verbose
}

//AddGoFile appends go file to slice gofiles
func (l *Output) AddGoFile(name string) {
	l.gofiles = append(l.gofiles, name)
}

//ProcessGoFiles processes go files
func (l *Output) ProcessGoFiles(work util.GenerationWork, noSkip bool) error {
	aggregates := make(map[string][]util.GeneratedCode)

	for generated := range work.Done {
		if generated.Error == util.ErrorSkip {
			Log(logSkipped+" Generation skipped [%s]", generated.Generator)
			l.skipped++
		} else if generated.Error != nil {
			Log(logError+" Generation failed [%s]: %s", generated.Generator, generated.Error)
			l.failure++
		} else if generated.Aggregate {
			a := aggregates[generated.Filename]
			aggregates[generated.Filename] = append(a, generated)
		} else {
			fname, lg, err := saveGenerated(generated, noSkip)
			Log(lg)

			if err == nil {
				if strings.HasSuffix(fname, ".go") {
					// AddGoFile(fname)
					l.AddGoFile(fname)
				} else if strings.HasSuffix(fname, ".sh") {
					os.Chmod(fname, 0755)
				}

				_written++
			} else if err == util.ErrorSkip {
				l.skipped++
			} else {
				l.failure++
			}
		}
		return nil
	}

	for _, generated := range aggregates {
		fname, lg, err := saveAggregate(generated, noSkip)
		Log(lg)

		if err == nil {
			if strings.HasSuffix(fname, ".go") {

				l.AddGoFile(fname)
			}

			_written++
		} else if err == util.ErrorSkip {
			l.skipped++
		} else {
			l.failure++
		}
	}

	return nil
}

// GenerateAndOverwrite deletes a file if exists and then saves it
func (l *Output) GenerateAndOverwrite(component string, template string, filename string, data interface{}) {
	var err error
	filename, err = util.GetAbsPath(filename)

	if err != nil {
		l.WithFields(log.Fields{"filename": filename, "error": err}).Error("An error occurred.")
		l.failure++
	} else if util.FileExists(filename) {
		l.WithFields(log.Fields{"filename": filename}).Warn("Deleting existing file.")
	}

	l.GenerateAndSave(component, template, filename, data)
}

// GenerateAndSave saves a generated file and returns error
func (l *Output) GenerateAndSave(component string, template string, filename string, data interface{}) {
	var (
		code string
		err  error
		mode os.FileMode = 0755
	)

	filename, err = util.GetAbsPath(filename)
	if err != nil {
		l.WithFields(log.Fields{"filename": filename, "error": err}).Error("An error occurred.")
		l.failure++
		return
	}

	if util.FileExists(filename) {
		l.WithFields(log.Fields{"filename": filename}).Warn("Skipping existing file.")
		l.skipped++
		return
	}

	if err = os.MkdirAll(path.Dir(filename), mode); err != nil {
		l.WithFields(log.Fields{"filename": filename, "error": err}).Error("An error occurred.")
		l.failure++
		return
	}

	code, err = util.ExecuteTemplate(template, data)
	if err != nil {
		l.WithFields(log.Fields{"filename": filename, "error": err}).Error("An error occurred.")
		l.failure++
		return
	}

	err = ioutil.WriteFile(filename, []byte(code), mode)
	if err != nil {
		l.WithFields(log.Fields{"filename": filename, "error": err}).Error("An error occurred.")
		l.failure++
		return
	}

	l.WithFields(log.Fields{"filename": filename}).Info("File written.")
	l.success++

	if strings.HasSuffix(filename, ".go") {
		l.AddGoFile(filename)
	} else if strings.HasSuffix(filename, ".sh") {
		os.Chmod(filename, 0755)
	}
}

// initToolset check if all required tools are present
func initToolset() {
	var (
		err error
		ok  = true
	)

	_tools.GoImports, err = exec.LookPath("goimports")
	if err != nil {
		fmt.Println("Required tool goimports not found: ", err)
		ok = false
	}

	_tools.GoFmt, err = exec.LookPath("gofmt")
	if err != nil {
		fmt.Println("Required tool gofmt not found: ", err)
		ok = false
	}

	_tools.Protoc, err = exec.LookPath("protoc")
	if err != nil {
		fmt.Println("Required tool protoc not found: ", err)
		ok = false
	}

	_, err = exec.LookPath("protoc-gen-go")
	if err != nil {
		fmt.Println("Required tool protoc-gen-go not found: ", err)
		fmt.Println("Install using go get -u github.com/golang/protobuf/protoc-gen-go")
		ok = false
	}

	_tools.Dep, err = exec.LookPath("dep")
	if err != nil {
		fmt.Println("Required tool dep not found: ", err)
		fmt.Println("Install using go get -u github.com/golang/dep/cmd/dep")
		ok = false
	}

	if !ok {
		log.Fatalln("Please install above tools before continuing.")
	}
}
