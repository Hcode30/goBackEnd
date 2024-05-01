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
  _,err := MakeUser("hcode19@gmail.com","Hazem","Ahmed","Hcode@123")
  if err != nil {
    log.Fatal(err)
  }
	server := myServer.Init()
	mux := myServer.mux
	mux.HandleFunc("/view/", MyHandler(ViewHandler))
	mux.HandleFunc("GET /signup/", signUpPageHandler)
	mux.HandleFunc("POST /signup", signUpHandler)
	mux.HandleFunc("GET /login/", logInPageHandler)
	mux.HandleFunc("POST /login", logInHandler)
	mux.HandleFunc("/viewAll/",ViewAllHandler)
	mux.HandleFunc("/edit/", MyHandler(EditHandler))
	mux.HandleFunc("POST /save/", MyHandler(SaveHandler))
	periodicalSave()
	PrepareShutDown(server)
	log.Fatal(myServer.ListenAndServe())
}
