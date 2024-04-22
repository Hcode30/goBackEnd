package main

import "net/http"


type Server struct {
	mux    *http.ServeMux
	port   string
	server *http.Server
}

func (s *Server) Init() *http.Server {
	s.server = &http.Server{Addr: s.port, Handler: s.mux}
	return s.server
}

func (s *Server) ListenAndServe() error {
  println("listening on port", s.port[1:])
	return s.server.ListenAndServe()
}


