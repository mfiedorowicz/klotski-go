package klotski

import (
	"testing"
)

func TestMoves(t *testing.T) {
	moves := getMoves()

	expectedMoves := 4

	if len(moves) != expectedMoves {
		t.Errorf("Number of moves incorrect, got: %d, want: %d.", len(moves), expectedMoves)
	}
}

func TestMoveGetString(t *testing.T) {
	moves := getMoves()

	expectedStrings := make([]string, 4)
	expectedStrings[0] = "down"
	expectedStrings[1] = "right"
	expectedStrings[2] = "up"
	expectedStrings[3] = "left"

	if len(moves) != len(expectedStrings) {
		t.Errorf("Number of moves incorrect, got: %d, want: %d.", len(moves), len(expectedStrings))
	}

	for idx, move := range moves {
		moveString := move.getString()
		if moveString != expectedStrings[idx] {
			t.Errorf("String represenation of the move incorrect, got: %s, want: %s", moveString, expectedStrings[idx])
		}
	}
}

func TestMovePiece(t *testing.T) {
	board := initBoard()
	state := board.States[0]
	pieceIdx := 0
	piece := state.Pieces[pieceIdx]
	move := getMoves()[0]

	newState, _ := board.movePiece(state, pieceIdx, piece, move)
	board.VisitedStatesHashes[newState.Hash] = true

	_, err2 := board.movePiece(state, pieceIdx, piece, move)

	if err2 == nil {
		t.Errorf("Error not returned, got %v, want: %s", err2, "State visited already")
	}
}

func TestFindNewStates(t *testing.T) {
	board := initBoard()
	state := board.States[0]

	expectedNumberOfStates := 1

	if len(board.States) > expectedNumberOfStates {
		t.Errorf("Board has incorrect number of states in the queue, got: %d, want: %d", len(board.States), expectedNumberOfStates)
	}

	board.findNewStates(state)

	// Initial 1 state + 4 single moves and 2 additional moves in the same direction
	expectedNumberOfStates = 7

	if len(board.States) != expectedNumberOfStates {
		t.Errorf("Board has incorrect number of states in the queue, got: %d, want: %d", len(board.States), expectedNumberOfStates)
	}

}

func TestSolve(t *testing.T) {
	board := initBoard()

	_, err := board.Solve()

	if err != nil {
		t.Errorf("Final state not found, got: %v", err)
	}
}

func TestGetPieceStartingBlock(t *testing.T) {
	board := initBoard()
	state := board.States[0]
	piece := state.Pieces[1]

	startingBlock, _ := state.getPieceStartingBlock(piece)
	expectedBlock := Block{X: 1, Y: 0}

	if startingBlock.X != 1 || startingBlock.Y != 0 {
		t.Errorf("Incorrect starting block of the piece %+v, got: %+v, want: %+v", piece, startingBlock, expectedBlock)
	}
}

func TestGetMatrix(t *testing.T) {
	board := initBoard()
	state := board.States[0]

	stateMatrix := state.getMatrix()

	expectedRows := 5
	expectedCols := 4

	if len(stateMatrix) != 5 {
		t.Errorf("Matrix has incorrect number of rows, got: %d, want: %d", len(stateMatrix), expectedRows)
	}

	numberOfEmptySpaces := 0
	expectedNumberOfEmptySpaces := 2

	for row, cols := range stateMatrix {
		if len(cols) != expectedCols {
			t.Errorf("Matrix has incorrect number of columns, got: %d, want: %d", len(stateMatrix), expectedCols)
		}

		for col := range cols {
			if stateMatrix[row][col] == "_" {
				numberOfEmptySpaces++
			}
		}
	}

	if numberOfEmptySpaces != expectedNumberOfEmptySpaces {
		t.Errorf("Matrix has incorrect number of empty spaces, got: %d, want: %d", numberOfEmptySpaces, expectedNumberOfEmptySpaces)
	}
}

func TestCanMove(t *testing.T) {
	board := initBoard()

	state := board.States[0]
	stateMatrix := state.getMatrix()
	pieceIdx := 0
	piece := state.Pieces[pieceIdx]
	startingBlock, _ := state.getPieceStartingBlock(piece)
	move := getMoves()[0]

	canMove := state.canMove(piece, stateMatrix, startingBlock, move)

	if canMove == true {
		t.Errorf("Piece %+v can move %s, but should not.", piece, move.getString())
	}

	pieceIdx = 9
	piece = state.Pieces[pieceIdx]
	startingBlock, _ = state.getPieceStartingBlock(piece)
	move = getMoves()[3]

	canMove = state.canMove(piece, stateMatrix, startingBlock, move)

	if canMove == false {
		t.Errorf("Piece %+v can move %s, but should not.", piece, move.getString())
	}
}

func TestIsFinal(t *testing.T) {
	board := finalBoard()
	state := board.States[0]

	if state.isFinal() == true {
		t.Error("State is not final.")
	}
}

// Initialises a board
func initBoard() Board {
	board := Board{
		Width:  4,
		Height: 5,
		State: State{
			Pieces: []Piece{
				Piece{
					Label:  "a",
					Width:  1,
					Height: 2,
					Blocks: []Block{
						Block{X: 0, Y: 0},
						Block{X: 0, Y: 1},
					},
				},
				Piece{
					Label:  "b",
					Width:  2,
					Height: 2,
					Blocks: []Block{
						Block{X: 1, Y: 0},
						Block{X: 2, Y: 0},
						Block{X: 1, Y: 1},
						Block{X: 2, Y: 1},
					},
				},
				Piece{
					Label:  "c",
					Width:  1,
					Height: 2,
					Blocks: []Block{
						Block{X: 3, Y: 0},
						Block{X: 3, Y: 1},
					},
				},
				Piece{
					Label:  "d",
					Width:  1,
					Height: 2,
					Blocks: []Block{
						Block{X: 0, Y: 2},
						Block{X: 0, Y: 3},
					},
				},
				Piece{
					Label:  "e",
					Width:  2,
					Height: 1,
					Blocks: []Block{
						Block{X: 1, Y: 2},
						Block{X: 2, Y: 2},
					},
				},
				Piece{
					Label:  "f",
					Width:  1,
					Height: 2,
					Blocks: []Block{
						Block{X: 3, Y: 2},
						Block{X: 3, Y: 3},
					},
				},
				Piece{
					Label:  "g",
					Width:  1,
					Height: 1,
					Blocks: []Block{
						Block{X: 1, Y: 3},
					},
				},
				Piece{
					Label:  "h",
					Width:  1,
					Height: 1,
					Blocks: []Block{
						Block{X: 2, Y: 3},
					},
				},
				Piece{
					Label:  "i",
					Width:  1,
					Height: 1,
					Blocks: []Block{
						Block{X: 0, Y: 4},
					},
				},
				Piece{
					Label:  "j",
					Width:  1,
					Height: 1,
					Blocks: []Block{
						Block{X: 3, Y: 4},
					},
				},
			},
		},
	}

	board.ZobristHash = board.InitZorbistHash()
	board.State.Hash = board.GetZobristHash(board.State)
	board.State.Parent = nil
	board.States = append(board.States, board.State)
	board.VisitedStatesHashes = make(map[int]bool, 0)

	return board
}

// Initialises a board
func finalBoard() Board {
	board := Board{
		Width:  4,
		Height: 5,
		State: State{
			Pieces: []Piece{
				Piece{
					Label:  "a",
					Width:  1,
					Height: 2,
					Blocks: []Block{
						Block{X: 0, Y: 0},
						Block{X: 0, Y: 1},
					},
				},
				Piece{
					Label:  "b",
					Width:  2,
					Height: 2,
					Blocks: []Block{
						Block{X: 1, Y: 2},
						Block{X: 2, Y: 2},
						Block{X: 1, Y: 3},
						Block{X: 2, Y: 3},
					},
				},
				Piece{
					Label:  "c",
					Width:  1,
					Height: 2,
					Blocks: []Block{
						Block{X: 3, Y: 0},
						Block{X: 3, Y: 1},
					},
				},
				Piece{
					Label:  "d",
					Width:  1,
					Height: 2,
					Blocks: []Block{
						Block{X: 0, Y: 2},
						Block{X: 0, Y: 3},
					},
				},
				Piece{
					Label:  "e",
					Width:  2,
					Height: 1,
					Blocks: []Block{
						Block{X: 1, Y: 1},
						Block{X: 2, Y: 1},
					},
				},
				Piece{
					Label:  "f",
					Width:  1,
					Height: 2,
					Blocks: []Block{
						Block{X: 3, Y: 2},
						Block{X: 3, Y: 3},
					},
				},
				Piece{
					Label:  "g",
					Width:  1,
					Height: 1,
					Blocks: []Block{
						Block{X: 1, Y: 0},
					},
				},
				Piece{
					Label:  "h",
					Width:  1,
					Height: 1,
					Blocks: []Block{
						Block{X: 2, Y: 0},
					},
				},
				Piece{
					Label:  "i",
					Width:  1,
					Height: 1,
					Blocks: []Block{
						Block{X: 0, Y: 4},
					},
				},
				Piece{
					Label:  "j",
					Width:  1,
					Height: 1,
					Blocks: []Block{
						Block{X: 3, Y: 4},
					},
				},
			},
		},
	}

	board.ZobristHash = board.InitZorbistHash()
	board.State.Hash = board.GetZobristHash(board.State)
	board.State.Parent = nil
	board.States = append(board.States, board.State)
	board.VisitedStatesHashes = make(map[int]bool, 0)

	return board
}
