package util

import (
	"io/ioutil"
	"os"
	"path"
)

// GenerateAndSave saves a generated file and returns error
func GenerateAndSave(template string, filename string, data interface{}, noOverwrite bool, noSkip bool) error {
	var (
		code     string
		err      error
		isString bool
		mode     os.FileMode = 0755
	)

	filename, err = GetAbsPath(filename)
	if err != nil {
		return err
	}

	if !noSkip && noOverwrite && FileExists(filename) {
		return ErrorSkip
	}

	if err = os.MkdirAll(path.Dir(filename), mode); err != nil {
		return err
	}

	if code, isString = data.(string); !isString {
		code, err = ExecuteTemplate(template, data)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(filename, []byte(code), mode)
	if err != nil {
		return err
	}

	return nil
}
