
# CHESS STUDY

---

The purpose of this program is for me to delve deeper into chess programming while also improving at chess.

My first big milestone is analyzing all of my available games from the Chess.com API to see the common positions 
in my games (FEN representation). This will help me improve as I progress in the London System, with the white pieces,
and a Scandanavian variation (1. e4 d5 2. exd5 Nf6), with the black pieces. 

I have set up an architecture in Go to hit various Chess.com endpoints (namely `https://api.chess.com/pub/player/{username}/games/archives`
and `https://api.chess.com/pub/player/{username}/games/{YYYY}/{MM}`), process the data, and send it to a PostgreSQL database hosted in a Docker 
container. Currently, I have succesffuly processed all of my 3979 games to 116244 unique board positions where it is my move.

Have made a lot of progress since I last updated the README. `compute.py` is a python script to break down pgns and turn them into a group of fens. Using `exec.Command()` was the
easiest method I found to call a python script in Go, and it suit my needs well. I successfully filled the games and bridge tables. There were errors related to illegal moves 
due to Chess960 games and games from a Hikaru viewer tournament where the starting position was not `rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`.

I then tried my luck with docker and docker-compose, to little avail. The mixed language Dockerfile for the API was not working well, and the docker-compose was also a bit hectic. 
I plan to work on classifying the different board postions, most likely using Go. This, I plan, will be a microservice architecture, which will get me more experience with 
docker-compose and Nginx. The immediate next steps, however, are finishing setting up the different necessary routes and working on a frontend. The only routes currently get 
relevant table counts, get the months and years of processed data, and update the games table with any new games. I also may remove the halfmove clock from the fens, as that is not
particularly important to this project.

Although I did work on some speed improvements, my 48 months, 3979 games, and around 150,000 rows of other data takes about 7 minutes to process from empty tables. 
I tried to process Hikaru's games... it took about 1hr15min to process 50,000 games, with over 1million board positions and over 2million entries in the bridge table. I am 
going to look into how this process can be sped up, because if this were to be a product for the public, the current processing would be severly underwhelming and unacceptable.

Two other considerations for future features/updates/improvements: an in-depth look at goroutines, to see how concurrent processing can help my situation and 
integrating testing, which I regretfully have not been integrating as I go. I will certainly get much more familiar with Go interfaces.
