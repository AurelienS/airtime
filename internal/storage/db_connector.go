package storage

import (
	"database/sql"
	"log"
	"os"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/AurelienS/airtime/internal/storage/ent"
	_ "github.com/jackc/pgx/v5/stdlib" // for postgres driver
)

var (
	connectionString = os.Getenv("DB")
)

func Open() *ent.Client {

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}
