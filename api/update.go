package api

import (
	"chess-study/chesscom"
	"chess-study/dockpg"
	"encoding/json"
	"net/http"
)

// FUNCITON FOR UPDATING ARCHIVES INDIVIDUALLY, NOT YET IMPLEMENTED
func (db *Database) UpdateArchives(w http.ResponseWriter, r *http.Request) {
	Archive := chesscom.PingArchive()

	Archive.Endpoints = dockpg.AddNewArchives(db.DB, Archive)

	newGames := dockpg.AddNewGames(db.DB, Archive.Endpoints)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newGames)
}

// FUNCTION FOR UPDATING GAMES INDIVIDUALLY, NOT YET IMPLEMENTED
func (db *Database) UpdateGames(w http.ResponseWriter, r *http.Request) {
	Archive := chesscom.PingArchive()

	Archive.Endpoints = dockpg.AddNewArchives(db.DB, Archive)

	newGames := dockpg.AddNewGames(db.DB, Archive.Endpoints)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newGames)
}

// UPDATE TABLES WITH ARCHIVES, GAMES, AND PROCESS GAME DATA FOR COUNTS AND BRIDGE
func (db *Database) Update(w http.ResponseWriter, r *http.Request) {
	Archive := chesscom.PingArchive()

	Archive.Endpoints = dockpg.AddNewArchives(db.DB, Archive)

	newGames := dockpg.AddNewGames(db.DB, Archive.Endpoints)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newGames)
}
