package dockpg

import (
	"chess-study/chesscom"
	"chess-study/helpers"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
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
var initial string

type AddedGames struct {
	Games []ToPython `json:"games"`
}

type AddedGame struct {
	Url   string `json:"url"`
	Pgn   string `json:"pgn"`
	Color bool   `json:"color"`
}

type ToPython struct {
	Pgn   string `json:"pgn"`
	Color bool   `json:"color"`
	// Url string `json:"url"`
}

type Counts struct {
	Counts []FenCount `json:"counts"`
}

type FenCount struct {
	Fen   string `json:"fen"`
	Count int    `json:"count"`
}

func AddNewGames(db *sql.DB, months []string) int {
	defer helpers.Duration(helpers.Track("AddNewGame"))

	addedGamesCount := 0
	lastGame := getLastGame(db)

	for monthEntry := range months {
		games := chesscom.PingMonth(months[monthEntry]).Games
		// fmt.Printf("url: %s\nNum of games: %d\n\n", months[monthEntry], len(games))
		currentIndex := len(games) - 1
		currentGame := games[currentIndex]

		for currentGame.URL != lastGame {
			url = currentGame.URL
			pgn = currentGame.PGN
			timecontrol = currentGame.TimeControl
			fen = currentGame.FEN
			color = determineColor(currentGame)
			initial = currentGame.Initial

			if color {
				myrating = currentGame.White.Rating
				opprating = currentGame.Black.Rating
				myresult = currentGame.White.Result
				oppresult = currentGame.Black.Result
				accuracy = currentGame.Accuracies.White
			} else {
				myrating = currentGame.Black.Rating
				opprating = currentGame.White.Rating
				myresult = currentGame.Black.Result
				oppresult = currentGame.White.Result
				accuracy = currentGame.Accuracies.Black
			}

			if initial == "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" {
				game := insertGame(db, url, pgn, timecontrol, myrating, opprating, color, fen, myresult, oppresult, accuracy)

				if game.Url != "" && game.Pgn != "" {
					addedGamesCount++

					args, err := json.Marshal(ToPython{Pgn: pgn, Color: color})

					cmd := exec.Command("python3", "fenpy/compute.py", string(args))
					out, err := cmd.CombinedOutput()
					if err != nil {
						log.Fatalf("Command execution failed: %s\nOutput: %s", err, out)
					}

					// fmt.Println("Output: ", string(out))

					var responseObj Counts
					if err := json.Unmarshal(out, &responseObj); err != nil {
						log.Fatalf("Failed to unmarshal JSON: %s\nJSON: %s", err, out)
					}
					pythonToFenAndBridge(db, responseObj, url)
				}
			}
			currentIndex = currentIndex - 1
			if currentIndex < 0 {
				break
			}
			currentGame = games[currentIndex]
		}
	}
	return addedGamesCount

}

func determineColor(game chesscom.Game) bool {
	return game.White.Username == os.Getenv("USERNAME")
}

func pythonToFenAndBridge(db *sql.DB, counts Counts, url string) {
	for _, fens := range counts.Counts {
		countQuery := `INSERT INTO counts (fen, count)
		VALUES ($1, $2)
		ON CONFLICT (fen) DO UPDATE SET count = counts.count + 1`

		bridgeQuery := `INSERT INTO bridge (link, fen)
		VALUES ($1, $2)`

		_, err := db.Exec(countQuery, fens.Fen, fens.Count)
		if err != nil {
			fmt.Println("COUNTS TABLE ERROR")
			log.Fatal(err)
		}

		_, err = db.Exec(bridgeQuery, url, fens.Fen)
		if err != nil {
			fmt.Println("BRIDGE TABLE ERROR")
			log.Fatal(err)
		}
	}
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

	return AddedGame{Url: ourl, Pgn: opgn, Color: ocolor}
}
