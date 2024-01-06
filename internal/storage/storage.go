package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/AurelienS/cigare/internal/storage/sqlc"
	_ "github.com/lib/pq" // add support for postgres
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "tennis"
	dbname   = "cigare"
)

func Open() (*sqlc.Queries, error) {
	url := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	conn, err := sql.Open("postgres", url)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return sqlc.New(conn), nil
}
