package main

import (
	rice "github.com/GeertJohan/go.rice"
	"github.com/fluxynet/gocipe/cmd"
	"github.com/fluxynet/gocipe/util"
)

//go:generate rice embed-go

// Versioning info
var (
	appVersion = "n/a"
	appCommit  = "n/a"
	appBuilt   = "n/a"
)

func init() {
	util.SetTemplates(rice.MustFindBox("templates"))
	cmd.SetVersionInfo(appVersion, appCommit, appBuilt)
}

func main() {
	cmd.Execute()
}
