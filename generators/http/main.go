package http

var tmplMain = `
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

//GenerateMain generates code to kick start the http server
func GenerateMain() (string, error) {
	return tmplMain, nil
}
