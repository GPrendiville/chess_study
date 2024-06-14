package main

import (
	"chess-study/api"
	"database/sql"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	// dockpg.Connection()

	connect := os.Getenv("POSTGRES_CONNECT")

	db, err := sql.Open("postgres", connect)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(os.Getenv("PORT"), db)
	if err := server.InitalizeAPI(); err != nil {
		log.Fatal(err)
	}
}
