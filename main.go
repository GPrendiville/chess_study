package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

type EndpointArchive struct {
	Endpoints []string `json:"archives"`
}

type GamesFromMonth struct {
	Games []Game `json:"games"`
}

type Game struct {
	URL         string   `json:"url"`
	PGN         string   `json:"pgn"`
	Accuracies  Accuracy `json:"accuracies"`
	FEN         string   `json:"fen"`
	TimeControl string   `json:"time_class"`
	White       Player   `json:"white"`
	Black       Player   `json:"black"`
}

type Accuracy struct {
	White float64 `json:"white"`
	Black float64 `json:"black"`
}

type Player struct {
	Rating   int    `json:"rating"`
	Result   string `json:"result"`
	Username string `json:"username"`
}

func main() {
	connect := os.Getenv("POSTGRES_CONNECT")

	db, err := sql.Open("postgres", connect)

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	response, err := http.Get(os.Getenv("CHESS_COM"))
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject EndpointArchive
	json.Unmarshal(responseData, &responseObject)

	responseObject.Endpoints = addNewArchives(db, responseObject)

	// month, err := http.Get(responseObject.Endpoints[0])
	// if err != nil {
	// 	fmt.Print(err.Error())
	// 	os.Exit(1)
	// }

	// monthData, err := io.ReadAll(month.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var monthObject GamesFromMonth
	// json.Unmarshal(monthData, &monthObject)

	// fmt.Println(monthObject)

	// for i := 0; i < len(monthObject.Games); i++ {
	// 	fmt.Println(monthObject.Games[i].Accuracies)
	// }
}

func addNewArchives(db *sql.DB, archives EndpointArchive) []string {
	var id int

	if checkArchivesData(db) {
		lastEndpointQuery := `SELECT id
		FROM archives
		ORDER BY id DESC
		LIMIT 1`

		err := db.QueryRow(lastEndpointQuery).Scan(&id)
		if err != nil {
			fmt.Println("GET LAST ENDPOINT ERROR")
			log.Fatal(err)
		}
	}

	for i := id; i < len(archives.Endpoints); i++ {
		insertEndpoint(db, archives.Endpoints[i])
	}

	if id > 0 {
		id--
	}

	return archives.Endpoints[id:]
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

func checkArchivesData(db *sql.DB) bool {
	var check bool
	query := `SELECT EXISTS(SELECT 1 FROM archives)`

	err := db.QueryRow(query).Scan(&check)
	if err != nil {
		fmt.Println("CHECK ARCHIVES DATA ERROR")
		log.Fatal(err)
	}

	return check
}
