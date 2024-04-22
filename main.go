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
	mux.HandleFunc("/view/", MyHandler(ViewHandler))
	mux.HandleFunc("/edit/", MyHandler(EditHandler))
	mux.HandleFunc("/save/", MyHandler(SaveHandler))
	PrepareShutDown(server)
	log.Fatal(myServer.ListenAndServe())
}
