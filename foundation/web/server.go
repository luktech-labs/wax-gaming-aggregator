package web

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type Server struct {
	router *mux.Router
	mws    []Middleware
}

func NewServer(mw ...Middleware) *Server {
	return &Server{router: mux.NewRouter(), mws: mw}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Subrouter(prefix string, mw ...Middleware) *Server {
	router := s.router.PathPrefix(prefix).Subrouter()
	mws := append(s.mws, mw...)

	return &Server{router: router, mws: mws}
}

func (s *Server) Handle(method, path string, handler Handler, mw ...Middleware) {
	// first wrap server specific middlewares
	h := wrapMiddleware(s.mws, handler)

	// wrap handler specific middlewares
	h = wrapMiddleware(mw, h)

	handleFn := func(w http.ResponseWriter, r *http.Request) {
		err := h(r.Context(), w, r)
		if err != nil {
			RespondError(r.Context(), w, err, http.StatusInternalServerError)
		}
	}

	s.router.HandleFunc(path, handleFn).Methods(method)
}
