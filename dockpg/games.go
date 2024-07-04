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
	"sync"
)

// GAME VARIABLES
var url string
var pgn string
var timecontrol string
var myrating int
var opprating int
var color bool
var fen string
var rules string
var myresult string
var oppresult string
var accuracy float64

// PROCESS ALL OUTSTANDING GAMES FROM EACH ARCHIVE AND GENERATE/INSERT COUNTS
func AddNewGames(db *sql.DB, months []string) int {
	defer helpers.Duration(helpers.Track("AddNewGame"))

	var gamesForFenGen AddedGames
	var addedGamesCount int

	results := make(chan AddedGame)
	done := make(chan bool)
	ch := make(chan ToPython, 10)

	var wg sync.WaitGroup

	go func() {
		for result := range results {
			gamesForFenGen.Games = append(gamesForFenGen.Games, result)
			addedGamesCount++
		}
		done <- true
	}()

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go gameRoutines(db, results, ch, &wg)
	}

	for monthEntry := range months {
		fmt.Printf("Archives: %s", months[monthEntry])
		games := chesscom.PingMonth(months[monthEntry]).Games

		for _, game := range games {

			if game.Initial == "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" && game.RULES == "chess" {
				url = game.URL
				pgn = game.PGN
				rules = game.RULES
				timecontrol = game.TimeControl
				fen = game.FEN
				color = determineColor(game)

				if color {
					myrating = game.White.Rating
					opprating = game.Black.Rating
					myresult = game.White.Result
					oppresult = game.Black.Result
					accuracy = game.Accuracies.White
				} else {
					myrating = game.Black.Rating
					opprating = game.White.Rating
					myresult = game.Black.Result
					oppresult = game.White.Result
					accuracy = game.Accuracies.Black
				}

				ch <- ToPython{url, pgn, timecontrol, myrating, opprating, rules, color, fen, myresult, oppresult, accuracy}
			}

		}
	}
	close(ch)
	wg.Wait()

	close(results)
	<-done

	if len(gamesForFenGen.Games) > 0 {

		args, _ := json.Marshal(gamesForFenGen)
		fmt.Println("Games marshalled")

		cmd := exec.Command("python3", "fenpy/compute.py")
		stdinPipe, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			defer stdinPipe.Close()
			io.WriteString(stdinPipe, string(args))
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

// DETERMINE PIECE COLOR OF PLAYER
func determineColor(game chesscom.Game) bool {
	return game.White.Username == os.Getenv("USERNAME")
}

// GO ROUTINES FOR INSERTING GAMES
func gameRoutines(db *sql.DB, results chan<- AddedGame, ch <-chan ToPython, wg *sync.WaitGroup) {
	defer wg.Done()

	for game := range ch {
		result := insertGame(db, game.Url, game.Pgn, game.TimeControl, game.MyRating, game.OppRating, game.Color, game.Fen, game.MyResult, game.OppResult, game.Accuracy)
		if result.Url != "" && result.Pgn != "" {
			results <- result
		}
	}
}

// INSERT GAME TO GAMES TABLE AND RETURN INFORMATION FOR FEN GENERATION IF SUCCESSFULL
func insertGame(db *sql.DB, url string, pgn string, timecontrol string, myrating int, opprating int, color bool, fen string, myresult string, oppresult string, accuracy float64) AddedGame {
	var ourl string
	var opgn string
	var ocolor bool
	query := `INSERT INTO games (url, pgn, time_control, my_rating, opp_rating, color, fen, my_result, opp_result, accuracy)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING url, pgn, color`

	row := db.QueryRow(query, url, pgn, timecontrol, myrating, opprating, color, fen, myresult, oppresult, accuracy)
	err := row.Scan(&ourl, &opgn, &ocolor)

	if err != nil {
		fmt.Println(err, " INSERT ERROR")
		return AddedGame{"", "", false}
	}

	return AddedGame{Url: ourl, Pgn: opgn, Color: ocolor}
}
