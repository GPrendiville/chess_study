package dockpg

import (
	"chess-study/chesscom"
	"database/sql"
	"fmt"
	"log"
	"os"
)

var url string
var pgn string
var timecontrol string
var myrating int
var opprating int
var color bool
var fen string
var myresult string
var oppresult string
var accuracy float64

type AddedGames struct {
	Games []AddedGame `json:"games"`
}

type AddedGame struct {
	Url   string `json:"url"`
	Pgn   string `json:"pgn"`
	Color bool   `json:"color"`
}

func AddNewGames(db *sql.DB, months []string) AddedGames {
	addedGames := AddedGames{Games: []AddedGame{}}
	lastGame := getLastGame(db)

	for monthEntry := range months {
		games := chesscom.PingMonth(months[monthEntry]).Games

		for i := 0; i < len(games); i++ {
			if games[i].URL != lastGame {
				url = games[i].URL
				pgn = games[i].PGN
				timecontrol = games[i].TimeControl
				fen = games[i].FEN
				color = determineColor(games[i])

				if color {
					myrating = games[i].White.Rating
					opprating = games[i].Black.Rating
					myresult = games[i].White.Result
					oppresult = games[i].Black.Result
					accuracy = games[i].Accuracies.White
				} else {
					myrating = games[i].Black.Rating
					opprating = games[i].White.Rating
					myresult = games[i].Black.Result
					oppresult = games[i].White.Result
					accuracy = games[i].Accuracies.Black
				}

				game := insertGame(db, url, pgn, timecontrol, myrating, opprating, color, fen, myresult, oppresult, accuracy)

				if game.Url != "" && game.Pgn != "" {
					addedGames.Games = append(addedGames.Games, game)
				}

			}
		}
	}
	return addedGames

}

func determineColor(game chesscom.Game) bool {
	return game.White.Username == os.Getenv("USERNAME")
}

func getLastGame(db *sql.DB) string {
	var url string

	if !checkGameTable(db) {
		return "NONE"
	}

	lastEndpointQuery := `SELECT url
		FROM games
		ORDER BY id DESC
		LIMIT 1`

	err := db.QueryRow(lastEndpointQuery).Scan(&url)
	if err != nil {
		fmt.Println("GET LAST GAME ERROR")
		log.Fatal(err)
	}

	return url
}

func checkGameTable(db *sql.DB) bool {
	var check bool
	query := `SELECT EXISTS(SELECT 1 FROM games)`

	err := db.QueryRow(query).Scan(&check)
	if err != nil {
		fmt.Println("CHECK GAMES TABLE ERROR")
		log.Fatal(err)
	}

	return check
}

func insertGame(db *sql.DB, url string, pgn string, timecontrol string, myrating int, opprating int, color bool, fen string, myresult string, oppresult string, accuracy float64) AddedGame {
	var ourl string
	var opgn string
	var ocolor bool
	query := `INSERT INTO games (url, pgn, time_control, my_rating, opp_rating, color, fen, my_result, opp_result, accuracy)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING url, pgn, color`

	row := db.QueryRow(query, url, pgn, timecontrol, myrating, opprating, color, fen, myresult, oppresult, accuracy)
	err := row.Scan(&ourl, &opgn, &ocolor)

	if err != nil {
		return AddedGame{"", "", false}
	}

	fmt.Println(ourl)

	return AddedGame{ourl, opgn, ocolor}
}
