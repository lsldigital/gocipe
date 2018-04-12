
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
