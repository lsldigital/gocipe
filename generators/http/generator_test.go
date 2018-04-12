package http

import (
	"log"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	var g generator

	g.Filename = ""
	g.FlagSet = nil
	output, err := g.Generate()
	log.Fatal(output)

	expected := `
package main
go:generate gocipe rest -file $GOFILE -struct Card
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
`

	if err != nil {
		t.Errorf("Got error: %s", err)
	} else if strings.Compare(output, expected) != 0 {
		t.Errorf(`Output not as expected. Output length = %d Expected length = %d
--- # Output ---
%s
--- # Expected ---
%s
------------------
`, len(output), len(expected), output, expected)
	}
}

// go:generate gocipe rest -file $GOFILE -struct Card
