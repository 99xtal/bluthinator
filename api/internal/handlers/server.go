package handlers

import (
	"database/sql"

	"github.com/99xtal/bluthinator/api/internal/services"
	"github.com/elastic/go-elasticsearch/v7"
)

type Server struct {
	DB            *sql.DB
	ES            *elasticsearch.Client
	ObjectStorage *services.StorageClient
}

func NewServer(db *sql.DB, es *elasticsearch.Client, storage *services.StorageClient) *Server {
	return &Server{DB: db, ES: es, ObjectStorage: storage}
}
