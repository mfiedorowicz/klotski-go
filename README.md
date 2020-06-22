# Klotski Go

## About

The purpose of this application is to solve a sliding puzzle board game known as Klotski (or Huarong Dao and many other names), in as few moves as possble, using breadth-first search (BFS) algorithm, written in Go programming language. For tracking visited states I've used Zorbist hashing function(https://en.wikipedia.org/wiki/Zobrist_hashing) used mainly for solving 2-dimensional board games, i.e. chess and Go.

This application solves the sliding blocks puzzle in 90 moves. If there are 2 empty spaces in given direction, the given piece can move 1 or 2 spaces (but only in the same direction as previous move) and counts as 1 move. There are solutions for 81 moves but these ones move the given piece by 2 spaces but in different direction.

## Running the application

- `make test` - runs unit tests

- `make build` - builds the executable file

- `make run-cli` - runs the application in the CLI mode

- `make run-http` - runs the application in the HTTP server mode
