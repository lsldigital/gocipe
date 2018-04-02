package generators

import (
	"os"
	"path"
)

//GetAbsPath returns absolute path
func GetAbsPath(location string) (string, error) {
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
