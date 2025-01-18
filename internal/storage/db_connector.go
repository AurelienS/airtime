package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/AurelienS/airtime/internal/storage/ent"
	_ "github.com/jackc/pgx/v5/stdlib" // for postgres driver
)

var (
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbname   = os.Getenv("DB_NAME")
	dbhost   = os.Getenv("DB_HOST")
)

// User ID=%s;Password=%s;Host=%s;Port=%s;Database=%s.
func Open() *ent.Client {
	databaseURL := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		user,
		password,
		dbhost,
		port,
		dbname,
	)
	fmt.Println("file: db_connector.go ~ line 32 ~ funcOpen ~ databaseURL : ", databaseURL)
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}
