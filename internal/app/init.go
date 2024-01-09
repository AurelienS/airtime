package app

import (
	"encoding/gob"

	"github.com/AurelienS/cigare/internal/auth"
	"github.com/AurelienS/cigare/internal/log"
	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/webserver"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5"
)

func Initialize(isProd bool) (*webserver.Server, error) {
	store := configureSessionStore(isProd)
	queries, db, err := initializeDatabase()
	if err != nil {
		return nil, err
	}

	server := webserver.NewServer(*queries, db, store)
	return server, nil
}

func configureSessionStore(isProd bool) sessions.Store {
	store := auth.NewStore(isProd)
	auth.ConfigureGoth(store)
	gob.Register(storage.User{})
	return store
}

func initializeDatabase() (*storage.Queries, *pgx.Conn, error) {
	db, err := storage.Open()
	if err != nil {
		log.Fatal().Msg("Cannot open db")
		return nil, nil, err
	}
	queries := storage.New(db)
	return queries, db, nil
}
