package http

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/fluxynet/gocipe/generators"
)

func init() {
	generators.AddCommand("http", "Generate http server", factory)
}

type generator struct {
	FlagSet  *flag.FlagSet
	Filename string
}

func factory(args []string) (generators.Command, error) {
	var g generator
	flagset := flag.NewFlagSet("http", flag.ExitOnError)
	flagset.StringVar(&g.Filename, "file", "", "Filename of main.go")

	flagset.Parse(args)

	return g, nil
}

func (g generator) Generate() (string, error) {
	var generated []string
	generated = append(generated, `
package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

`)

	if segment, err := GenerateContainer(); err == nil {
		generated = append(generated, segment)
	} else {
		return "", err
	}

	if segment, err := GenerateMain(); err == nil {
		generated = append(generated, segment)
	} else {
		return "", err
	}

	targetFilename := filepath.Dir(g.Filename) + "/http.go"
	output := strings.Join(generated, "\n")
	err := ioutil.WriteFile(targetFilename, []byte(output), 0644)

	return output, err
}
