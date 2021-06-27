package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPServer interface {
	Router() http.Handler
	GetHealth(w http.ResponseWriter, r *http.Request)
	Serve() error
}

type httpServer struct {
	router   http.Handler
	httpAddr string
}

func New(addr string) HTTPServer {
	server := &httpServer{
		httpAddr: addr,
	}
	router(server)

	return server
}

func router(s *httpServer) {
	r := mux.NewRouter()

	r.HandleFunc("/health", s.GetHealth).Methods(http.MethodGet)

	s.router = r
}

func (s *httpServer) Router() http.Handler {
	return s.router
}

func (s *httpServer) GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "ok")
}

func (s *httpServer) Serve() error {
	return http.ListenAndServe(s.httpAddr, s.Router())
}
