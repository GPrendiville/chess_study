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

func AddNewGames(db *sql.DB, months []string) {

	for monthEntry := range months {
		games := chesscom.PingMonth(months[monthEntry]).Games

		for i := 0; i < len(games); i++ {
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

			insertGame(db, url, pgn, timecontrol, myrating, opprating, color, fen, myresult, oppresult, accuracy)
		}
	}

}

func determineColor(game chesscom.Game) bool {
	return game.White.Username == os.Getenv("USERNAME")
}

func insertGame(db *sql.DB, url string, pgn string, timecontrol string, myrating int, opprating int, color bool, fen string, myresult string, oppresult string, accuracy float64) {
	query := `INSERT INTO games (url, pgn, time_control, my_rating, opp_rating, color, fen, my_result, opp_result, accuracy)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (url) 
		DO NOTHING`

	result, err := db.Exec(query, url, pgn, timecontrol, myrating, opprating, color, fen, myresult, oppresult, accuracy)
	if err != nil {
		fmt.Println("INSERT ARCHIVE ERROR")
		log.Fatal(err)
	}
	_ = result
}
