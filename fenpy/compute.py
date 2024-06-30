#!/usr/local/bin/python
import json
import re
from dataclasses import dataclass
import chess
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

    return [moves_arrs, list_of_moves, pgn_cleaned]

def removeMoveClocks(fen):
    elements = fen.split()
    return ' '.join(elements[:4])

def gen_fen(game, color):
    board = chess.Board()
    group = []

    def return_fen():
        return color == (board.ply() % 2 == 0)

    for move in game[0]:
        try:
            board.push_san(move.white)
            if return_fen() and board.ply() > 1:
                group.append(removeMoveClocks(board.fen()))
        except chess.IllegalMoveError:
            print(f"Illegal move found for white: {move.white}")

        if move.black != '':
            try:
                board.push_san(move.black)
                if return_fen():
                    group.append(removeMoveClocks(board.fen()))
            except chess.IllegalMoveError:
                print(f"Illegal move found for black: {move.black}")

    board.reset()
    return group


def counting_fens(data):
    fen_data = {}

    for game in data['games']:
        pgn = game['pgn']
        color = game['color']
        urls = [game['url']]
        fens = gen_fen(extract_moves(pgn), color)

        for fen in fens:
            if fen in fen_data:
                # Update count and append the URL if it's not already included
                fen_data[fen]['count'] += 1
                if game['url'] not in fen_data[fen]['urls']:
                    fen_data[fen]['urls'].append(game['url'])
            else:
                # Initialize for new FEN entries
                fen_data[fen] = {'fen': fen, 'count': 1, 'urls': urls}

    # Prepare the output format
    counts_list = [{'fen': key, 'count': value['count'], 'urls': value['urls']} for key, value in fen_data.items()]
    return {'counts': counts_list}
    

r =  sys.stdin.read()

data = json.loads(r)

out = counting_fens(data)

print(json.dumps(out))