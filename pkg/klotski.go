package klotski

import (
	"bytes"
	"errors"
	"math"
	"math/rand"
	"sort"
	"time"
)

// Board defines a board, stores an initial state and all states leading to the final state.
// Also keep track of visited states.
type Board struct {
	Width               int
	Height              int
	ZobristHash         [][][]int
	State               State
	States              []State
	VisitedStatesHashes map[int]bool
}

// State defines a state of the board for each move, has reference to a parent state (before the move).
type State struct {
	Pieces        []Piece
	Hash          int
	Parent        *State
	Step          int
	MoveDirection string
	MovePiece     Piece
}

// Piece holds informatation about a piece on a board. Each piece is a combination of one or more single blocks.
type Piece struct {
	Label  string
	Width  int
	Height int
	Blocks []Block
}

// Block represents coordinates of a piece.
type Block struct {
	X, Y int
}

// Move defines a movement of a piece / block.
type Move struct {
	X, Y int
}

// Returns a list of possible moves.
func getMoves() []Move {
	moves := make([]Move, 4, 4)

	moves[0] = Move{0, 1}  // DOWN
	moves[1] = Move{1, 0}  // RIGHT
	moves[2] = Move{0, -1} // UP
	moves[3] = Move{-1, 0} // LEFT

	return moves
}

// Returns string represenation of a move.
func (m *Move) getString() string {
	var moveString string

	if m.X == 0 && m.Y > 0 {
		moveString = "down"
	} else if m.X == 0 && m.Y < 0 {
		moveString = "up"
	} else if m.X > 0 && m.Y == 0 {
		moveString = "right"
	} else if m.X < 0 && m.Y == 0 {
		moveString = "left"
	}

	return moveString
}

// InitZorbistHash initialises Zorbist hash for the board.
// Reference: https://en.wikipedia.org/wiki/Zobrist_hashing
func (board *Board) InitZorbistHash() [][][]int {
	const rows int = 5
	const cols int = 4

	zobristTable := make([][][]int, rows)
	for row := 0; row < rows; row++ {
		zobristTable[row] = make([][]int, cols)

		for col := 0; col < cols; col++ {
			zobristTable[row][col] = make([]int, 5)
		}
	}

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			for idx := 0; idx < 5; idx++ {
				rand.Seed(time.Now().UTC().UnixNano())
				zobristTable[row][col][idx] = int(rand.Float64() * math.Pow(2.0, 31.0))
			}
		}
	}

	return zobristTable
}

// GetZobristHash returns hash of a state.
func (board *Board) GetZobristHash(state State) int {
	const rows int = 5
	const cols int = 4

	hash := 0

	var pieceTypes [2][2]int

	pieceTypes[0][1] = 1
	pieceTypes[1][1] = 2
	pieceTypes[1][0] = 3
	pieceTypes[0][0] = 4

	for _, piece := range state.Pieces {
		for _, block := range piece.Blocks {
			x, y := block.X, block.Y
			pieceType := pieceTypes[piece.Width-1][piece.Height-1]
			hash ^= board.ZobristHash[y][x][pieceType]
		}
	}

	return hash
}

// Returns updated Zorbist hash for a moved piece.
func (board *Board) getUpdatedZobristHash(state State, piece Piece, move Move) int {
	hash := state.Hash

	var pieceTypes [2][2]int

	pieceTypes[0][1] = 1
	pieceTypes[1][1] = 2
	pieceTypes[1][0] = 3
	pieceTypes[0][0] = 4

	pieceType := pieceTypes[piece.Width-1][piece.Height-1]

	for _, block := range piece.Blocks {
		hash ^= board.ZobristHash[block.Y][block.X][pieceType]
		hash ^= board.ZobristHash[block.Y][block.X][0]
		hash ^= board.ZobristHash[block.Y+move.Y][block.X+move.X][0]
		hash ^= board.ZobristHash[block.Y+move.Y][block.X+move.X][pieceType]
	}

	return hash
}

// Moves a piece in a given direction.
// Returns a new state or error if state been visited already.
func (board *Board) movePiece(state State, pieceIdx int, piece Piece, move Move) (State, error) {

	hash := board.getUpdatedZobristHash(state, piece, move)

	_, visited := board.VisitedStatesHashes[hash]

	if visited {
		return State{}, errors.New("State visited already")
	}

	var movedBlocks []Block
	for _, block := range piece.Blocks {
		newBlock := Block{X: block.X + move.X, Y: block.Y + move.Y}
		movedBlocks = append(movedBlocks, newBlock)
	}

	newPiece := Piece{
		Label:  piece.Label,
		Width:  piece.Width,
		Height: piece.Height,
		Blocks: movedBlocks,
	}

	newPieces := make([]Piece, len(state.Pieces))

	copy(newPieces, state.Pieces)

	newPieces[pieceIdx] = newPiece

	newState := State{
		Pieces:        newPieces,
		Parent:        &state,
		Step:          state.Step + 1,
		MovePiece:     newPiece,
		MoveDirection: move.getString(),
		Hash:          state.Hash,
	}

	newState.Hash = hash

	return newState, nil
}

// Finds new states for all possible (and not visited) moves and adds them the board states.
func (board *Board) findNewStates(state State) {

	stateMatrix := state.getMatrix()

	for pieceIdx, piece := range state.Pieces {

		startingBlock, _ := state.getPieceStartingBlock(piece)

		for _, move := range getMoves() {
			canMove := state.canMove(piece, stateMatrix, startingBlock, move)

			if canMove {
				newState, err := board.movePiece(state, pieceIdx, piece, move)
				if err == nil {
					_, visited := board.VisitedStatesHashes[newState.Hash]
					if visited == false {
						board.VisitedStatesHashes[newState.Hash] = true
						board.States = append(board.States, newState)

						newStateMatrix := newState.getMatrix()
						startingBlockOnNewState, _ := newState.getPieceStartingBlock(piece)
						pieceInNewState := newState.Pieces[pieceIdx]
						canMoveAgain := newState.canMove(pieceInNewState, newStateMatrix, startingBlockOnNewState, move)

						if canMoveAgain {
							newState2, err2 := board.movePiece(newState, pieceIdx, pieceInNewState, move)

							if err2 == nil {
								_, visited := board.VisitedStatesHashes[newState2.Hash]

								if visited == false {
									board.VisitedStatesHashes[newState2.Hash] = true
									newState2.Step--
									board.States = append(board.States, newState2)
								}
							}
						}
					}
				}
			}
		}
	}
}

// Solve finds a solution for the initial board state
func (board *Board) Solve() ([]State, error) {

	results := make([]State, 0)

	for idx := 0; idx < len(board.States); idx++ {

		currentState := board.States[idx]

		board.VisitedStatesHashes[currentState.Hash] = true

		if currentState.isFinal() {

			end := false
			newState := currentState

			uniqueStates := make(map[int]State, 0)

			idx2 := 1

			for end == false {

				if newState.Parent != nil || newState.Step > 0 {
					uniqueStates[newState.Step] = newState

					idx2++

					if newState.MovePiece.Label == newState.Parent.MovePiece.Label && newState.MoveDirection == newState.Parent.MoveDirection {
						newState = *newState.Parent
					}

					newState = *newState.Parent
				} else {
					end = true
				}
			}

			var steps []int

			for step := range uniqueStates {
				steps = append(steps, step)
			}

			sort.Ints(steps)

			for _, step := range steps {
				results = append(results, uniqueStates[step])
			}

			return results, nil
		}

		board.findNewStates(currentState)
	}

	return results, errors.New("Cannot solve")
}

// Returns starting block (top left one) of a piece.
func (state *State) getPieceStartingBlock(piece Piece) (Block, error) {
	const rows int = 5
	const cols int = 4

	var startingBlock Block

	for _, p := range state.Pieces {
		x, y := cols, rows
		for _, b := range p.Blocks {
			if b.X < x {
				x = b.X
			}

			if b.Y < y {
				y = b.Y
			}
		}

		if p.Label == piece.Label {
			startingBlock.X = x
			startingBlock.Y = y

			return startingBlock, nil
		}
	}

	return startingBlock, errors.New("Cannot find piece starting block")
}

// Returns a matrix for a given state.
func (state *State) getMatrix() [][]string {
	const rows int = 5
	const cols int = 4

	boardMatrix := make([][]string, rows)

	for row := 0; row < rows; row++ {
		boardMatrix[row] = make([]string, cols)
	}

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			boardMatrix[row][col] = "_"
		}
	}

	for _, piece := range state.Pieces {
		for _, block := range piece.Blocks {
			boardMatrix[block.Y][block.X] = piece.Label
		}
	}

	return boardMatrix
}

// Checks if a piece can be moved in a given direction.
func (state *State) canMove(piece Piece, boardMatrix [][]string, startingBlock Block, move Move) bool {
	const rows int = 5
	const cols int = 4

	canMove := false

	switch move.getString() {
	case "down":
		for x := startingBlock.X; x < (startingBlock.X + piece.Width); x++ {
			y := startingBlock.Y + piece.Height
			canMove = y <= rows-1 && x <= cols-1 && boardMatrix[y][x] == "_"

			if canMove == false {
				break
			}
		}
	case "right":
		for y := startingBlock.Y; y < (startingBlock.Y + piece.Height); y++ {
			x := startingBlock.X + piece.Width
			canMove = y <= rows-1 && x <= cols-1 && boardMatrix[y][x] == "_"

			if canMove == false {
				break
			}
		}
	case "up":
		for x := startingBlock.X; x < (startingBlock.X + piece.Width); x++ {
			y := startingBlock.Y - 1
			canMove = y >= 0 && x <= cols-1 && boardMatrix[y][x] == "_"

			if canMove == false {
				break
			}
		}
	case "left":
		for y := startingBlock.Y; y < (startingBlock.Y + piece.Height); y++ {
			x := startingBlock.X - 1
			canMove = y <= rows-1 && x >= 0 && boardMatrix[y][x] == "_"

			if canMove == false {
				break
			}
		}
	}

	return canMove
}

// Checks if a state is a final one.
func (state *State) isFinal() bool {
	for _, piece := range state.Pieces {
		if piece.Label == "b" {
			startingBlock, _ := state.getPieceStartingBlock(piece)

			if startingBlock.Y == 3 && startingBlock.X == 1 {
				return true
			}
		}
	}

	return false
}

// Print a given board state
func (board *Board) Print(state State) string {
	const rows int = 5
	const cols int = 4

	var buffer bytes.Buffer

	var stateMatrix [rows][cols]string

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			stateMatrix[row][col] = "_"
		}
	}

	for _, piece := range state.Pieces {
		for _, block := range piece.Blocks {
			stateMatrix[block.Y][block.X] = piece.Label
		}
	}

	for rowIdx := 0; rowIdx < board.Width+2; rowIdx++ {
		buffer.WriteString("X ")
	}

	buffer.WriteString("\n")

	for rowIdx := 0; rowIdx < board.Height; rowIdx++ {
		buffer.WriteString("X ")
		for colIdx := 0; colIdx < board.Width; colIdx++ {
			buffer.WriteString(stateMatrix[rowIdx][colIdx] + " ")
		}
		buffer.WriteString("X \n")
	}
	for rowIdx := 0; rowIdx < board.Width+2; rowIdx++ {
		if rowIdx == 2 || rowIdx == 3 {
			buffer.WriteString("Z ")
		} else {
			buffer.WriteString("X ")
		}
	}

	buffer.WriteString("\n")

	return buffer.String()
}
