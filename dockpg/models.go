package dockpg

// MODELS FOR DOCK PACKAGE

type AddedGames struct {
	Games []AddedGame `json:"games"`
}

type AddedGame struct {
	Url   string `json:"url"`
	Pgn   string `json:"pgn"`
	Color bool   `json:"color"`
}

type ToPython struct {
	Url         string
	Pgn         string
	TimeControl string
	MyRating    int
	OppRating   int
	Rules       string
	Color       bool
	Fen         string
	MyResult    string
	OppResult   string
	Accuracy    float64
}

type Counts struct {
	Counts []FenCount `json:"counts"`
}

type FenCount struct {
	Fen   string   `json:"fen"`
	Count int      `json:"count"`
	Urls  []string `json:"urls"`
}
