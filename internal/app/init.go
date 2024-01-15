package app

import (
	"context"
	"log"

	"github.com/AurelienS/cigare/internal/storage"
	"github.com/AurelienS/cigare/internal/storage/ent"
	"github.com/AurelienS/cigare/internal/storage/ent/migrate"
	"github.com/AurelienS/cigare/web"
	"github.com/AurelienS/cigare/web/session"
)

func Initialize(isProd bool) (*web.Server, error) {
	store := session.ConfigureSessionStore(isProd)
	client := initializeDatabase()

	if err := client.Debug().Schema.Create(
		context.Background(),
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	server := web.NewServer(client, store)
	return server, nil
}

func initializeDatabase() *ent.Client {
	return storage.Open()
}
