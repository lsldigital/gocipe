package output

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/fluxynet/gocipe/util"
	log "github.com/sirupsen/logrus"
)

type Output struct {
	*log.Logger
	gofiles                   []string
	success, failure, skipped int
}

func (l *Output) AddGoFile(name string) {
	l.gofiles = append(l.gofiles, name)
}

func ProcessGoFiles() error {
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
