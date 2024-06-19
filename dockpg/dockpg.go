package dockpg

import (
	"database/sql"
	"log"
	"os"
)

func Connection() *sql.DB {
	connect := os.Getenv("POSTGRES_CONNECT")

	db, err := sql.Open("postgres", connect)

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
