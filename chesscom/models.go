package chesscom

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
