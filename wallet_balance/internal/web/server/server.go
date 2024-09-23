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

func (s *Server) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *Server) Start() {
	s.Router.Use(middleware.Logger)
	for path, handler := range s.Handlers {
		s.Router.Post(path, handler)
	}

	http.ListenAndServe(s.ServerPort, s.Router)
}
