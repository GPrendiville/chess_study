#!/usr/local/bin/python
import json
import re
from dataclasses import dataclass
import chess
import requests
import os
from dotenv import load_dotenv
from collections import Counter
import sys


@dataclass
class Move():
    white: str
    black: str = ""


@dataclass
class Game():
    pgn: str
    color: bool


def string_to_move(move_string):
    try:
        white, black = move_string.strip().split()
        return Move(white, black)
    except:
        return Move(move_string.strip())


def extract_moves(pgn):
    pgn_body = pgn.split('\n\n', 1)[1]

    pgn_no_comments = re.sub(r'\{.*?\}', '', pgn_body).rsplit(' ', 1)[0]

    moves = re.sub(r' \d+\.\.\. ', '', pgn_no_comments)

    pgn_cleaned = ' '.join(moves.split())

    list_of_moves = re.split(r'\d+\.', pgn_cleaned)[1:]

    moves_arrs = [string_to_move(move) for move in list_of_moves]

    return moves_arrs

def gen_fen(game, color):
    board = chess.Board()
    group = []

    def return_fen():
        return color == (board.ply() % 2 != 0)

    for move in game:
        try:
            board.push_san(move.white)
            if return_fen() and board.ply() > 1:
                group.append(board.fen().split(' ', 1)[0])
        except chess.IllegalMoveError:
            print(f"Illegal move found: {move.white}")

        if move.black != '':
            try:
                board.push_san(move.black)
                if return_fen():
                    group.append(board.fen().split(' ', 1)[0])
            except chess.IllegalMoveError:
                print(f"Illegal move found: {move.black}")

    board.reset()
    return group


def counting_fens(data):
    games = []
    counts = Counter()
    counted = {'fencounts': []}

    # moves = extract_moves(pgn)
    # fens = gen_fen(moves, False)

    ####### NEEEEEEEEEEED RULES TO AVOID CHESS960

    for game in data['games']:
        if game['rules'] != 'chess':
            continue
        games.append(Game(game['pgn'], (game['white']['username'] == 'C4DDY903')))

    for game in games:
        counts.update(gen_fen(extract_moves(game.pgn), game.color))

    for position in counts:
        if (counts[position] > 4):
            counted['fencounts'].append({'fen': position, 'count': counts[position]})

    return counted

r =  sys.argv[1]
# header = {'User-Agent': 'griffithprendiville@gmail.com'}
# requests.get("https://api.chess.com/pub/player/c4ddy903/games/2024/06", headers=header).json()
# r.json()
data = json.loads(r)

out = counting_fens(data)

print(json.dumps(out))