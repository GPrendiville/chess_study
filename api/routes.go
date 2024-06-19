package api

import (
	"chess-study/chesscom"
	"encoding/json"
	"log"
	"net/http"

	"chess-study/dockpg"
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

func (db *Database) UpdateArchives(w http.ResponseWriter, r *http.Request) {
	Archive := chesscom.PingArchive()

	Archive.Endpoints = dockpg.AddNewArchives(db.DB, Archive)

	newGames := dockpg.AddNewGames(db.DB, Archive.Endpoints)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newGames)
}

func (db *Database) UpdateGames(w http.ResponseWriter, r *http.Request) {
	Archive := chesscom.PingArchive()

	Archive.Endpoints = dockpg.AddNewArchives(db.DB, Archive)

	newGames := dockpg.AddNewGames(db.DB, Archive.Endpoints)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newGames)
}

func (db *Database) Update(w http.ResponseWriter, r *http.Request) {
	Archive := chesscom.PingArchive()

	Archive.Endpoints = dockpg.AddNewArchives(db.DB, Archive)

	newGames := dockpg.AddNewGames(db.DB, Archive.Endpoints)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newGames)
}
