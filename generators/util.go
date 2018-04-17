package generators

import (
	"os"
	"path"
	"strings"
)

//GetAbsPath returns absolute path
func GetAbsPath(location string) (string, error) {
	location = strings.TrimRight(location, "\\/")
	if path.IsAbs(location) {
		return location, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	location = path.Clean(wd + "/" + location)
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
