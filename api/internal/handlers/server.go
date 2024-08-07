package handlers

import (
	"database/sql"

	"github.com/elastic/go-elasticsearch/v7"
)

type Server struct {
	DB *sql.DB
	ES *elasticsearch.Client
}

func NewServer(db *sql.DB, es *elasticsearch.Client) *Server {
	return &Server{DB: db, ES: es}
}