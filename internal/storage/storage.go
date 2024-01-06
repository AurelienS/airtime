package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // add support for postgres
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "tennis"
	dbname   = "cigare"
)

func Open() (*Queries, error) {
	url := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	conn, err := sql.Open("postgres", url)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return New(conn), nil
}
