package seeder

import (
	"fmt"
	"log"
	"os"

	"github.com/lsldigital/gocipe/util"
)

// Generate returns generated database schema creation code
func Generate(r *util.Recipe) {

	filename, err := util.GetAbsPath("schema/seeder.gocipe.sql")
	if err != nil {
		log.Fatal("Cannot get absolute path", err)
	}
	statements := util.GenerataSeeds(r)
	fi, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Cannot create f", err)
	}
	defer fi.Close()

	for _, s := range statements {
		fmt.Fprintf(fi, s)
	}
}
