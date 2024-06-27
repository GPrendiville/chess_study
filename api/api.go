package api

import (
	"database/sql"
	"net/http"
)

type Database struct {
	*sql.DB
}

type APIServer struct {
	APIEndpoint string
	db          *sql.DB
}

func NewServer(APIEndpoint string, db *sql.DB) *APIServer {
	return &APIServer{
		APIEndpoint: APIEndpoint,
		db:          db,
	}
}

func (api *APIServer) InitalizeAPI() error {
	router := http.NewServeMux()
	db := &Database{
		api.db,
	}
	router.HandleFunc("GET /", db.GetCounts)
	router.HandleFunc("GET /archives", db.GetArchives)
	router.HandleFunc("GET /update", db.Update)

	server := http.Server{
		Addr:    api.APIEndpoint,
		Handler: router,
	}

	return server.ListenAndServe()
}
