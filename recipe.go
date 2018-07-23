package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fluxynet/gocipe/output"
	"github.com/fluxynet/gocipe/util"
)

func loadRecipe() (*util.Recipe, error) {
	var recipe util.Recipe

	recipePath, err := util.GetAbsPath("gocipe.json")
	if err != nil {
		return nil, err
	}

	if !util.FileExists(recipePath) {
		return nil, fmt.Errorf("file not found: %s", recipePath)
	}

	recipeContent, err := ioutil.ReadFile(recipePath)

	if err != nil {
		return nil, err
	}

	output.Log("%x", sha256.Sum256([]byte(recipeContent)))

	output.Inject(
		recipePath,
	)

	err = json.Unmarshal(recipeContent, &recipe)
	if err != nil {
		return nil, fmt.Errorf("recipe decoding failed: %s", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	util.WorkingDir = wd
	util.ProjectImportPath = strings.TrimPrefix(wd, os.Getenv("GOPATH")+"/src/")
	os.Getenv("GOPATH")

	return &recipe, nil
}
