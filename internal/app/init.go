package app

import (
	"encoding/gob"
	"log"

	"github.com/AurelienS/cigare/internal/auth"
	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/webserver"
	"github.com/gorilla/sessions"
)

func Initialize(isProd bool) *webserver.Server {
	store := configureSessionStore(isProd)
	queries := initializeDatabase()

	server := webserver.NewServer(*queries, store)
	return server
}

func configureSessionStore(isProd bool) sessions.Store {
	store := auth.NewStore(isProd)
	auth.ConfigureGoth(store)
	gob.Register(storage.User{})
	return store
}

func initializeDatabase() *storage.Queries {
	queries, err := storage.Open()
	if err != nil {
		log.Fatal("Cannot open db")
		return nil
	}
	return queries
}
