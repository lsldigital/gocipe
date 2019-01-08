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

	//Delete file if already exists
	exists(filename)

	fi, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Cannot create f", err)
	}
	defer fi.Close()

	for _, s := range statements {
		fmt.Fprintf(fi, s)
	}
}

func exists(filename string) {

	_, err := os.Stat(filename)

	if !os.IsNotExist(err) {
		err1 := os.Remove(filename)
		if err1 != nil {
			log.Println("File cannot be deleted: " + filename)
		}
	}
}
