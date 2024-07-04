
# CHESS STUDY

---

The purpose of this program is for me to delve deeper into chess programming while also improving at chess.

My first big milestone is analyzing all of my available games from the Chess.com API to see the common positions 
in my games (FEN representation). This will help me improve as I progress in the London System, with the white pieces,
and a Scandanavian variation (1. e4 d5 2. exd5 Nf6), with the black pieces. 

I have set up an architecture in Go to hit various Chess.com endpoints (namely `https://api.chess.com/pub/player/{username}/games/archives`
and `https://api.chess.com/pub/player/{username}/games/{YYYY}/{MM}`), process the data, and send it to a PostgreSQL database hosted in a Docker 
container.

`compute.py` is a python script to break down pgns and turn them into a group of fens. Using `exec.Command()` was theeasiest method I found to call a 
python script in Go, and it suit my needs well. I successfully filled the games and bridge tables. There were errors related to illegal moves 
due to Chess960 games and games from a Hikaru viewer tournament where the starting position was not `rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`.

I have succesffuly imporved the speed of the `update` endpoint. My 4000 games went from about 6 min 50 sec to 40 sec. I also tested with GM Hikaru Nakamura's Chess.com
games, going from 1hr15min to around 5 min 15 sec. To do this, I utilized Go routines for inserting games and batching for inserting to the `counts` and `bridge` tables.
I plan to also optimize the python script, as that is where most of the time spent is at this iteration of the API.

The next major tickets to tackle:
1. Python script optimization
2. Go testing
3. FEN endpoint and determine best construction of information
4. Frontend development

Notes on #3:
    The averages of the counts for both my 4000 games and Hikaru's 50000+ games are very similar. The overall average is around 1.1 and the average
    when `count` is above 10 hovers around 46.0. This leads me to believe I can construct the query to be a generic solution to providing the right
    amount of information.
