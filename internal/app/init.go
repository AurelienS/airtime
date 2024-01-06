package app

import (
	"encoding/gob"
	"log"

	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/webserver"
	"github.com/AurelienS/cigare/internal/webserver/session"
	"github.com/gorilla/sessions"
)

func Initialize(isProd bool) *webserver.Server {
	store := configureSessionStore(isProd)
	queries := initializeDatabase()

	server := webserver.NewServer(queries, store)
	return server
}

func configureSessionStore(isProd bool) sessions.Store {
	store := session.NewStore(isProd)
	session.ConfigureGoth(store)
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
