package main

import (
	"chess-study/api"
	"chess-study/dockpg"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	db := dockpg.Connection()

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	server := api.NewServer(os.Getenv("PORT"), db)
	if err := server.InitalizeAPI(); err != nil {
		log.Fatal(err)
	}
}
