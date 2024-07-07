package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Nandainthegrass/Flipzon/Authentication/services/user"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(router)

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
