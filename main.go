package main

import (
	"log"
	"net/http"
)

var myServer = &Server{
	mux:  http.NewServeMux(),
	port: ":8000",
}

func main() {
	server := myServer.Init()
	mux := myServer.mux
	// mux.HandleFunc("/view/", MyHandler(ViewHandler))
	mux.HandleFunc("/home", homePageHandler)
	mux.HandleFunc("GET /signup/", signUpPageHandler)
	mux.HandleFunc("POST /api/v1/signup", signUpHandler)
	mux.HandleFunc("GET /login/", logInPageHandler)
	mux.HandleFunc("POST /api/v1/login", logInHandler)
	// mux.HandleFunc("/viewAll/",ViewAllHandler)
	// mux.HandleFunc("/edit/", MyHandler(EditHandler))
	// mux.HandleFunc("POST /save/", MyHandler(SaveHandler))
	// periodicalSave()
	PrepareShutDown(server)
	log.Fatal(myServer.ListenAndServe())
}
