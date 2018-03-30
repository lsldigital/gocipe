package http

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

//Command the name of the command to start this generator
const Command = "http"

//Generator represent arguments accepted by this generator
type Generator struct {
	FlagSet  *flag.FlagSet
	Filename string
}

//Description the description of the command when used on cli
const Description = "Generate http server"

//NewGenerator returns a new Generator
func NewGenerator() *Generator {
	arguments := new(Generator)
	arguments.FlagSet = flag.NewFlagSet("http", flag.ExitOnError)
	arguments.FlagSet.StringVar(&arguments.Filename, "file", "", "Filename of main.go")

	return arguments
}

//Generate produce the generated code according to options
func Generate(generator Generator) string {
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
		log.Fatalf("An error occured during GenerateContainer: %s\n", err)
	}

	if segment, err := GenerateMain(); err == nil {
		generated = append(generated, segment)
	} else {
		log.Fatalf("An error occured during GenerateMain: %s\n", err)
	}

	targetFilename := filepath.Dir(generator.Filename) + "/http.go"
	output := strings.Join(generated, "\n")
	ioutil.WriteFile(targetFilename, []byte(output), 0644)

	return output
}
