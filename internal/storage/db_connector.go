package storage

import (
	"database/sql"
	"fmt"
	"log"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/AurelienS/cigare/internal/storage/ent"
	_ "github.com/jackc/pgx/v5/stdlib" // for postgres driver
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "tennis"
	dbname   = "cigare"
)

func Open() *ent.Client {
	databaseURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}
