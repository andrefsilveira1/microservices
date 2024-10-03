package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Router     chi.Router
	Handlers   map[string]http.HandlerFunc
	ServerPort string
}

func NewServer(serverPort string) *Server {
	return &Server{
		Router:     chi.NewRouter(),
		Handlers:   make(map[string]http.HandlerFunc),
		ServerPort: serverPort,
	}
}

func (s *Server) AddHandler(method, path string, handler http.HandlerFunc) {
	switch method {
	case http.MethodGet:
		s.Router.Get(path, handler)
	case http.MethodPost:
		s.Router.Post(path, handler)
	case http.MethodPut:
		s.Router.Put(path, handler)
	case http.MethodDelete:
		s.Router.Delete(path, handler)
	}
}

func (s *Server) Start() {
	s.Router.Use(middleware.Logger)
	http.ListenAndServe(s.ServerPort, s.Router)
}
