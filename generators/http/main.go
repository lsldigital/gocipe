package http

var tmplMain = `
func main() {
	var err error
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	router := mux.NewRouter()
	
	if err = bootstrap(router); err != nil {
		log.Fatal("Bootstrap failed: ", err)
	}
	
	go func() {
		err = http.ListenAndServe(":" + port, router)
		if err != nil {
			log.Fatal("Failed to start http server: ", err)
		}
	}()

	log.Println("Listening on :" + port)
	<-sigs
	log.Println("Server stopped")
}
`

//GenerateMain generates code to kick start the http server
func GenerateMain() (string, error) {
	return tmplMain, nil
}
