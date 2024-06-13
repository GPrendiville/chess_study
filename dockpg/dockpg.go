package dockpg

import (
	"database/sql"
	"log"
	"os"

	"chess-study/chesscom"
)

func Connection() {
	connect := os.Getenv("POSTGRES_CONNECT")

	db, err := sql.Open("postgres", connect)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	Archive := chesscom.PingArchive()

	Archive.Endpoints = addNewArchives(db, Archive)

	AddNewGames(db, Archive.Endpoints)

}
