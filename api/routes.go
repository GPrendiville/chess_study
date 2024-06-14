package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type TotalArchives struct {
	Archives []Archive `json:"archives"`
}

type Archive struct {
	Id       int    `json:"id"`
	Endpoint string `json:"endpoint"`
}

func (db *Database) GetArchives(w http.ResponseWriter, r *http.Request) {
	archives := TotalArchives{Archives: []Archive{}}
	// make([]Archive, 0)

	archivesQuery := `SELECT *
		FROM archives`

	rows, err := db.DB.Query(archivesQuery)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id int
		var endpoint string
		if err = rows.Scan(&id, &endpoint); err != nil {
			log.Fatal(err)
		}
		archives.Archives = append(archives.Archives, Archive{id, endpoint})
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(archives)

}
