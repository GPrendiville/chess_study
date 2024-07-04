package dockpg

import (
	"database/sql"
	"fmt"
	"log"

	"chess-study/chesscom"
)

// INSERT NEW ARCHIVES TO ARCHIVES TABLE AND RETURN ARCHIVE LIST INCLUDING LAST INSERTED ARCHIVE TO PROCESS ANY NEW GAMES
func AddNewArchives(db *sql.DB, archives chesscom.EndpointArchive) []string {
	var id int

	if checkArchiveTable(db) {
		id = lastArchiveEntry(db)
	}

	for i := id; i < len(archives.Endpoints); i++ {
		insertEndpoint(db, archives.Endpoints[i])
	}

	if id > 0 {
		id--
	}

	return archives.Endpoints[id:]
}

// GET THE MOST RECENT ARCHIVE ENTRY FOR UPDATED GAME ARCHIVE
func lastArchiveEntry(db *sql.DB) int {
	var id int

	lastEndpointQuery := `SELECT id
		FROM archives
		ORDER BY id DESC
		LIMIT 1`

	err := db.QueryRow(lastEndpointQuery).Scan(&id)
	if err != nil {
		fmt.Println("GET LAST ENDPOINT ERROR")
		log.Fatal(err)
	}

	return id
}

// INSERT ARCHIVE TO TABLE
func insertEndpoint(db *sql.DB, endpoint string) {
	query := `INSERT INTO archives (endpoint)
		VALUES ($1)`

	result, err := db.Exec(query, endpoint)
	if err != nil {
		fmt.Println("INSERT ARCHIVE ERROR")
		log.Fatal(err)
	}
	_ = result
}

// CHECK IF ARCHIVE TABLE IS EMPTY
func checkArchiveTable(db *sql.DB) bool {
	var check bool
	query := `SELECT EXISTS(SELECT 1 FROM archives)`

	err := db.QueryRow(query).Scan(&check)
	if err != nil {
		fmt.Println("CHECK ARCHIVES DATA ERROR")
		log.Fatal(err)
	}

	return check
}
