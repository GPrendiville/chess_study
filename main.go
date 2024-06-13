package main

import (
	"chess-study/dockpg"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	dockpg.Connection()
}
