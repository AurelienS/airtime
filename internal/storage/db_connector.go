package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/AurelienS/cigare/internal/storage/ent"
	_ "github.com/jackc/pgx/v5/stdlib" // for postgres driver
)

var (
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbname   = os.Getenv("DB_NAME")
)

func Open() *ent.Client {
	databaseURL := fmt.Sprintf("postgres://%s:%s@db:%s/%s",
		user,
		password,
		port,
		dbname,
	)
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}
