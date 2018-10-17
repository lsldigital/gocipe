package output

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/fluxynet/gocipe/util"
	log "github.com/sirupsen/logrus"
)

var (
	_recipePath string
	_tools      toolset
	_verbose    bool
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

// Inject gets path injected into this package
func Inject(path string) {
	_recipePath = path
}

// SetVerbose can be used to toggle verbosity
func SetVerbose(verbose bool) {
	_verbose = verbose
}

//AddGoFile appends go file to slice gofiles
func (l *Output) AddGoFile(name string) {
	l.gofiles = append(l.gofiles, name)
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

// PostProcessGoFiles executes goimports and gofmt on go files that have been generated
func (l *Output) PostProcessGoFiles() {
	if len(l.gofiles) == 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(l.gofiles))

	for _, file := range l.gofiles {
		go func(file string) {
			defer wg.Done()

			cmd := exec.Command(_tools.GoImports, "-w", file)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()

			if err != nil {
				fmt.Printf("Error running %s on %s: %s\n", _tools.GoImports, file, err)
				return
			}

			cmd = exec.Command(_tools.GoFmt, "-w", file)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()

			if err != nil {
				fmt.Printf("Error running %s on %s: %s\n", _tools.GoFmt, file, err)
			}
		}(file)
	}

	wg.Wait()

	var mode string
	if util.FileExists(util.WorkingDir + "/Gopkg.toml") {
		mode = "ensure"
	} else {
		mode = "init"
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		cmd := exec.Command(_tools.Dep, mode)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			fmt.Printf("Error running dep %s: %s\n", mode, err)
		}
	}()

	fmt.Printf("dep %s in progress...", mode)
	wg.Wait()
}

// ProcessProto executes protoc to generate go files from protobuf files
func (l *Output) ProcessProto() {
	var (
		cmd    *exec.Cmd
		err    error
		mode   os.FileMode = 0755
		gopath             = os.Getenv("GOPATH") + "/src/"
	)

	// l.Info(logInfo + " Executing protoc to generate go files...")
	// l.failure++
	l.WithFields(log.Fields{"filename": "Protoc", "info": "Executing protoc to generate go files..."}).Info("Executing protoc.")

	// models.proto
	if !util.FileExists(util.WorkingDir + "/models") {
		if err = os.MkdirAll(util.WorkingDir+"/models", mode); err != nil {

			l.WithFields(log.Fields{"filename": " could not create folder: " + util.WorkingDir + "/models", "error": err}).Error("An error occurred.")
			l.failure++
			fmt.Printf(logError+" Error creating folder %s: %s\n", util.WorkingDir+"/models", err)

			return
		}
		l.WithFields(log.Fields{"Created folder": util.WorkingDir + "/models", "info": "Success creating folder: " + util.WorkingDir + "/models"}).Info("Create folder")
	}
	cmd = exec.Command(
		_tools.Protoc,
		`-I=`+util.WorkingDir+`/proto`,
		util.WorkingDir+`/proto/models.proto`,
		`--go_out=plugins=grpc:`+gopath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	if err != nil {

		l.WithFields(log.Fields{"filename": " models.proto", "error": err}).Error("An error occurred.")
		l.failure++
		return
	}

	l.WithFields(log.Fields{"Created folder": util.WorkingDir + "/models", "info": "Success creating folder: " + util.WorkingDir + "/models"}).Info("Create folder")

	// service_admin.proto, if admin service is to be generated
	if util.FileExists(util.WorkingDir + `/proto/service_admin.proto`) {
		if !util.FileExists(util.WorkingDir + "/services/admin") {
			if err = os.MkdirAll(util.WorkingDir+"/services/admin", mode); err != nil {
				l.WithFields(log.Fields{"Create folder": util.WorkingDir + "/services/admin", "error": err}).Error("An error occurred.")
				return
			}
			l.WithFields(log.Fields{"Create folder": "Success creating folder: " + util.WorkingDir + "/models"}).Info("Create folder")
		}
		cmd = exec.Command(
			_tools.Protoc,
			`-I=`+util.WorkingDir+`/proto`,
			util.WorkingDir+`/proto/service_admin.proto`,
			`--go_out=plugins=grpc:`+gopath,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()

		if err != nil {
			// Log(logError+" protoc execution error (%s): %s", "service_admin.proto", err)
			return
		}

		l.WithFields(log.Fields{"Protoc": "service_admin.proto", "info": "protoc generated go files from service_admin.proto"}).Info("Generate Protoc")

	}

	// cmd = exec.Command(
	// 	toolset.Protoc,
	// 	`-I=proto`,
	// 	`--plugin="protoc-gen-ts=`+util.WorkingDir+`/web/node_modules/.bin/protoc-gen-ts"`,
	// 	`--js_out="binary:`+util.WorkingDir+`/web/src/services"`,
	// 	`--ts_out="`+util.WorkingDir+`/web/src/services"`,
	// 	`proto/models.proto`,
	// )

	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// cmd.Dir = util.WorkingDir
	// err = cmd.Run()

	// if err != nil {
	// 	fmt.Printf("Error running %s: %s\n", toolset.Protoc, err)
	// 	return
	// }
}
