package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Nandainthegrass/Flipzon/Authentication/services/user"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
	rdb  *redis.Client
}

func NewAPIServer(addr string, db *sql.DB, rdb *redis.Client) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
		rdb:  rdb,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore, s.rdb)
	userHandler.RegisterRoutes(router)

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
