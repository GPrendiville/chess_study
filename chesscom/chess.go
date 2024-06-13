package chesscom

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func PingArchive() EndpointArchive {
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

	return responseObject
}

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
