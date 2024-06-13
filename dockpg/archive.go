package dockpg

import (
	"database/sql"
	"fmt"
	"log"

	"chess-study/chesscom"
)

func addNewArchives(db *sql.DB, archives chesscom.EndpointArchive) []string {
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
