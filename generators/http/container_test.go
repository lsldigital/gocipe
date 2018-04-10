package http

import (
	"strings"
	"testing"
)

//GenerateContainer generates code to load configuration for the application
func TestGenerateContainer(t *testing.T) {
	output, err := GenerateContainer()

	expected := `
const (
	//EnvironmentProd represents production environment
	EnvironmentProd = "PROD"
	//EnvironmentDev represents development environment
	EnvironmentDev  = "DEV"
)

var (
	env string
	dsn string
	db  *sql.DB
)

func init() {
	var (
		err error
	)

	godotenv.Load()

	env = os.Getenv("ENV")
	dsn = os.Getenv("DSN")

	if env == "" {
		log.Fatal("Environment variable ENV must be defined. Possible values are: DEV PROD")
	}

	if dsn == "" {
		log.Fatal("Environment variable DSN must be defined. Example: postgres://user:pass@host/db?sslmode=disable")
	}

	db, err = sql.Open("postgres", dsn)
	if err == nil {
		log.Println("Connected to database successfully.")
	} else if (env == EnvironmentDev) {
		log.Println("Database connection failed: ", err)
	} else {
		log.Fatal("Database connection failed: ", err)
	}

	err = db.Ping()
	if err == nil {
		log.Println("Pinged database successfully.")
	} else if (env == EnvironmentDev) {
		log.Println("Database ping failed: ", err)
	} else {
		log.Fatal("Database ping failed: ", err)
	}

	inject()
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
