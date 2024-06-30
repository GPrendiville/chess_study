package dockpg

import (
	"chess-study/chesscom"
	"chess-study/helpers"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
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
	Games []AddedGame `json:"games"`
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
	Fen   string   `json:"fen"`
	Count int      `json:"count"`
	Urls  []string `json:"urls"`
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
		gamesForFenGen := AddedGames{[]AddedGame{}}

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

					gamesForFenGen.Games = append(gamesForFenGen.Games, game)

				}
			}
			currentIndex = currentIndex - 1
			if currentIndex < 0 {
				break
			}
			currentGame = games[currentIndex]
		}

		args, _ := json.Marshal(gamesForFenGen)

		cmd := exec.Command("python3", "fenpy/compute.py")
		stdinPipe, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			defer stdinPipe.Close()
			io.WriteString(stdinPipe, string(args)) // `data` is the byte slice
		}()

		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("Command execution failed: %s\nOutput: %s", err, out)
		}

		var responseObj Counts
		if err := json.Unmarshal(out, &responseObj); err != nil {
			log.Fatalf("Failed to unmarshal JSON: %s\nJSON: %s", err, out)
		}

		pythonToFenAndBridge(db, responseObj.Counts)

	}
	return addedGamesCount
}

func determineColor(game chesscom.Game) bool {
	return game.White.Username == os.Getenv("USERNAME")
}

func pythonToFenAndBridge(db *sql.DB, data []FenCount) {
	const maxParams = 65535

	var countValues []string
	var countArgs []interface{}
	var bridgeValues []string
	var bridgeArgs []interface{}
	countIdx, bridgeIdx := 1, 1

	for _, entry := range data {
		if len(countArgs)+len(bridgeArgs)+len(entry.Urls)*2+2 > maxParams {
			// Execute batch
			executeBatch(db, countValues, countArgs, bridgeValues, bridgeArgs)
			countValues, countArgs, bridgeValues, bridgeArgs = nil, nil, nil, nil
			countIdx, bridgeIdx = 1, 1
		}

		// Append data for counts
		countValues = append(countValues, fmt.Sprintf("($%d, $%d)", countIdx, countIdx+1))
		countArgs = append(countArgs, entry.Fen, entry.Count)
		countIdx += 2

		// Handle multiple URLs for each FEN
		for _, url := range entry.Urls {
			bridgeValues = append(bridgeValues, fmt.Sprintf("($%d, $%d)", bridgeIdx, bridgeIdx+1))
			bridgeArgs = append(bridgeArgs, url, entry.Fen)
			bridgeIdx += 2
		}
	}

	if len(countArgs) > 0 || len(bridgeArgs) > 0 {
		executeBatch(db, countValues, countArgs, bridgeValues, bridgeArgs)
	}
}

func executeBatch(db *sql.DB, countValues []string, countArgs []interface{}, bridgeValues []string, bridgeArgs []interface{}) {

	countQuery := fmt.Sprintf("INSERT INTO counts (fen, count) VALUES %s ON CONFLICT (fen) DO UPDATE SET count = counts.count + excluded.count", strings.Join(countValues, ", "))
	bridgeQuery := fmt.Sprintf("INSERT INTO bridge (link, fen) VALUES %s", strings.Join(bridgeValues, ", "))

	tx, err := db.Begin()
	if err != nil {
		fmt.Println("ON BEGIN")
		log.Fatal(err)
	}

	_, err = tx.Exec(countQuery, countArgs...)
	if err != nil {
		tx.Rollback()
		fmt.Println("ON COUNTS EXEC")
		log.Fatal(err)
	}

	_, err = tx.Exec(bridgeQuery, bridgeArgs...)
	if err != nil {
		tx.Rollback()
		fmt.Println("ON BRIDGE EXEC")
		log.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		fmt.Println("ON COMMIT EXEC")
		log.Fatal(err)
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
