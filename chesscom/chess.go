package chesscom

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// GET ALL AVAILABLE ARCHIVES OF PLAYER AND UNPACK INTO GOLANG MODEL
func PingArchive() EndpointArchive {
	url := "https://api.chess.com/pub/player/" + os.Getenv("USERNAME") + "/games/archives"
	response, err := http.Get(url)
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

	return responseObject
}

// GET GAMES FROM ARCHIVE ENDPOINT AND UNPACK INTO GOLANG MODEL
func PingMonth(endpoint string) GamesFromMonth {
	response, err := http.Get(endpoint)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject GamesFromMonth
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}
