package app

import (
	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/util"
	"github.com/AurelienS/cigare/web"
	"github.com/AurelienS/cigare/web/session"
	"github.com/jackc/pgx/v5"
)

func Initialize(isProd bool) (*web.Server, error) {
	store := session.ConfigureSessionStore(isProd)
	queries, db, err := initializeDatabase()
	if err != nil {
		return nil, err
	}

	server := web.NewServer(*queries, db, store)
	return server, nil
}

func initializeDatabase() (*storage.Queries, *pgx.Conn, error) {
	db, err := storage.Open()
	if err != nil {
		util.Fatal().Msg("Cannot open db")
		return nil, nil, err
	}
	queries := storage.New(db)
	return queries, db, nil
}
