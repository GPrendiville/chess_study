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

def gen_fen(game, color, pgn):
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
            print(f"Illegal move found for white: {move.white} (pgn: {pgn})")

        if move.black != '':
            try:
                board.push_san(move.black)
                if return_fen():
                    group.append(removeMoveClocks(board.fen()))
            except chess.IllegalMoveError:
                print(f"Illegal move found for black: {move.black} (pgn: {pgn})")

    board.reset()
    return group


def counting_fens(data):
    pgn = data['pgn']
    color = data['color']
    counts = Counter()
    counted = {'counts': []}

    # moves = extract_moves(pgn)
    # fens = gen_fen(moves, False)

    ####### NEEEEEEEEEEED RULES TO AVOID CHESS960

    counts.update(gen_fen(extract_moves(pgn), color, data['pgn']))

    for position in counts:
        counted['counts'].append({'fen': position, 'count': counts[position]})

    return counted

r =  sys.argv[1]

data = json.loads(r)

out = counting_fens(data)

print(json.dumps(out))