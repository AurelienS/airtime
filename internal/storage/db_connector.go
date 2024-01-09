package storage

import (
	"context"
	"fmt"

	"github.com/AurelienS/cigare/internal/util"
	"github.com/jackc/pgx/v5"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "tennis"
	dbname   = "cigare"
)

func Open() (*pgx.Conn, error) {
	ctx := context.Background()

	url := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	conn, err := pgx.Connect(ctx, url)

	if err != nil {
		util.Error().Msgf("Failed to connect to postgres %s", err)
		return nil, err
	}

	return conn, nil
}
