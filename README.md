
# CHESS STUDY

---

The purpose of this program is for me to delve deeper into chess programming while also improving at chess.

My first big milestone is analyzing all of my available games from the Chess.com API to see the common positions 
in my games (FEN representation). This will help me improve as I progress in the London System, with the white pieces,
and a Scandanavian variation (1. e4 d5 2. exd5 Nf6), with the black pieces.

So far, I have set up an architecture in Go to hit various Chess.com endpoints (namely `https://api.chess.com/pub/player/{username}/games/archives`
and `https://api.chess.com/pub/player/{username}/games/{YYYY}/{MM}`), process the data, and send it to a PostgreSQL database hosted in a Docker container.

The next steps, in no particular order, are:
1. Setup a Go API to access and update database data
2. Setup a Python program to analyze the PGNs and send FENs to Postgres via the Go API
3. Create a new table to store the FENs and another table to relate them to the games of origin

There will certainly be a large amount of data relating to the FENs. The theoretical maximum is every move I have played/witnessed divided by 2, as
I only care about board positions where it is my turn to move. As I get more acquanited with the issue, I will assess the validity of decreasing the volume.



























