package app

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

// ServeHttp starts an http server
func ServeHttp(listen string, route func(router *mux.Router) error) {
	var err error
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	router := mux.NewRouter()
	err = route(router)

	if err != nil {
		log.Fatal("Failed to register routes: ", err)
	}

	go func() {
		err = http.ListenAndServe(listen, router)
		if err != nil {
			log.Fatal("Failed to start http server: ", err)
		}
	}()

	log.Println("Listening on " + listen)
	<-sigs
	log.Println("Server stopped")
}
