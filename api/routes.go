package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type TotalArchives struct {
	Archives []Archive `json:"archives"`
}

type Archive struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

func (db *Database) GetCounts(w http.ResponseWriter, r *http.Request) {
	var months int
	var games int
	var positions int

	tablesQuery := `SELECT n_live_tup
	FROM pg_stat_user_tables;`

	rows, err := db.DB.Query(tablesQuery)
	if err != nil {
		log.Fatal(err)
	}

	rows.Next()
	if rows.Scan(&months); err != nil {
		log.Fatal(err)
	}
	rows.Next()
	if rows.Scan(&games); err != nil {
		log.Fatal(err)
	}
	rows.Next()
	if rows.Scan(&positions); err != nil {
		log.Fatal(err)
	}
	rows.Close()

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	testingSentence := fmt.Sprintf("Months: %d\nGames: %d\nPositions: %d", months, games, positions)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(testingSentence)

}

func (db *Database) GetArchives(w http.ResponseWriter, r *http.Request) {
	archives := TotalArchives{Archives: []Archive{}}

	archivesQuery := `SELECT endpoint
		FROM archives`

	rows, err := db.DB.Query(archivesQuery)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var endpoint string
		if err = rows.Scan(&endpoint); err != nil {
			log.Fatal(err)
		}
		endpointArray := strings.Split(endpoint, "/")
		endpointArrayLength := len(endpointArray)
		month, _ := strconv.Atoi(endpointArray[endpointArrayLength-1])
		year, _ := strconv.Atoi(endpointArray[endpointArrayLength-2])

		archives.Archives = append(archives.Archives, Archive{month, year})
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(archives)

}
