package http

import (
	"strings"
	"testing"
)

//GenerateContainer generates code to load configuration for the application
func TestGenerateMain(t *testing.T) {
	output, err := GenerateMain()

	expected := `
func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	router := mux.NewRouter()
	registerRoutes(router)
	
	go func() {
		err := http.ListenAndServe(":8888", router)
		if err != nil {
			log.Fatal("Failed to start http server: ", err)
		}
	}()

	log.Println("Listening on :8888")
	<-sigs
	log.Println("Server stopped")
}
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
