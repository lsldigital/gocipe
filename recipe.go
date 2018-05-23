package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fluxynet/gocipe/util"
)

func loadRecipe() (*util.Recipe, error) {
	var recipe util.Recipe

	if len(os.Args) == 1 {
		return nil, errors.New("no recipe provided. usage: gocipe gocipe.json")
	}

	recipePath, err := util.GetAbsPath(os.Args[len(os.Args)-1])
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

	_recipeHash = fmt.Sprintf("%x", sha256.Sum256([]byte(recipeContent)))
	_recipePath = recipePath

	err = json.Unmarshal(recipeContent, &recipe)
	if err != nil {
		return nil, fmt.Errorf("recipe decoding failed: %s", err)
	}

	return &recipe, nil
}
