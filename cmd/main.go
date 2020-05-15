package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	klotski "github.com/mfiedorowicz/klotski-go/pkg"
)

var (
	mode = flag.String("mode", "cli", "run mode (cli default)")
)

func main() {

	flag.Parse()

	switch *mode {
	case "cli":
		runCli()
	case "http":
		runServer()
	}
}

func runCli() {
	board := initBoard()
	initialState := board.State
	results, err := board.Solve()

	if err != nil {
		fmt.Printf("Error occured: %s", err)
	} else {
		fmt.Printf("\nInitial State:\n\n")
		fmt.Println(board.Print(initialState))

		fmt.Printf("\nNumber of moves needed to reach final state: %d\n\n", len(results))
		for step, state := range results {
			fmt.Printf("%d) %s moves %s\n\n", step+1, state.MovePiece.Label, state.MoveDirection)
			fmt.Println(board.Print(state))
		}
	}
}

func runServer() {
	router := mux.NewRouter()

	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/solution", solutionHTMLPage).Methods("GET")
	log.Println("Listening on port 8000. Open http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	board := initBoard()
	initialState := board.State

	results, err := board.Solve()

	var resultsHTML template.HTML

	if err == nil {
		var buffer bytes.Buffer

		buffer.WriteString(fmt.Sprintf("<p>Number of moves needed to reach final state: <strong>%d</strong></p>", len(results)))

		for step, state := range results {
			buffer.WriteString("<div class=\"state\">")
			buffer.WriteString(fmt.Sprintf("<p>%d) <strong>%s</strong> moves <strong>%s</strong></p>", step+1, state.MovePiece.Label, state.MoveDirection))
			buffer.WriteString(strings.Replace(board.Print(state), "\n", "<br>", -1))
			buffer.WriteString("</div>")
		}

		resultsHTML = template.HTML(buffer.String())
	}

	title := "Klotski Go"

	data := struct {
		Title        string
		InitialState template.HTML
		Solution     template.HTML
	}{
		title,
		template.HTML(strings.Replace(board.Print(initialState), "\n", "<br>", -1)),
		resultsHTML,
	}

	tpl := template.Must(template.ParseFiles("cmd/templates/layout.html"))
	tpl.Execute(w, data)
}

func solutionHTMLPage(w http.ResponseWriter, r *http.Request) {
	board := initBoard()
	results, err := board.Solve()

	if err == nil {
		var buffer bytes.Buffer

		buffer.WriteString(fmt.Sprintf("<p>Number of moves: %d</p>", len(results)))

		resultsHTML := template.HTML(buffer.String())

		title := "Klotski Go"

		data := struct {
			Title   string
			Results template.HTML
		}{
			title,
			resultsHTML,
		}

		tpl := template.Must(template.ParseFiles("cmd/templates/solution.html"))
		tpl.Execute(w, data)
	}
}

// Initialises a board
func initBoard() klotski.Board {
	board := klotski.Board{
		Width:  4,
		Height: 5,
		State: klotski.State{
			Pieces: []klotski.Piece{
				klotski.Piece{
					Label:  "a",
					Width:  1,
					Height: 2,
					Blocks: []klotski.Block{
						klotski.Block{X: 0, Y: 0},
						klotski.Block{X: 0, Y: 1},
					},
				},
				klotski.Piece{
					Label:  "b",
					Width:  2,
					Height: 2,
					Blocks: []klotski.Block{
						klotski.Block{X: 1, Y: 0},
						klotski.Block{X: 2, Y: 0},
						klotski.Block{X: 1, Y: 1},
						klotski.Block{X: 2, Y: 1},
					},
				},
				klotski.Piece{
					Label:  "c",
					Width:  1,
					Height: 2,
					Blocks: []klotski.Block{
						klotski.Block{X: 3, Y: 0},
						klotski.Block{X: 3, Y: 1},
					},
				},
				klotski.Piece{
					Label:  "d",
					Width:  1,
					Height: 2,
					Blocks: []klotski.Block{
						klotski.Block{X: 0, Y: 2},
						klotski.Block{X: 0, Y: 3},
					},
				},
				klotski.Piece{
					Label:  "e",
					Width:  2,
					Height: 1,
					Blocks: []klotski.Block{
						klotski.Block{X: 1, Y: 2},
						klotski.Block{X: 2, Y: 2},
					},
				},
				klotski.Piece{
					Label:  "f",
					Width:  1,
					Height: 2,
					Blocks: []klotski.Block{
						klotski.Block{X: 3, Y: 2},
						klotski.Block{X: 3, Y: 3},
					},
				},
				klotski.Piece{
					Label:  "g",
					Width:  1,
					Height: 1,
					Blocks: []klotski.Block{
						klotski.Block{X: 1, Y: 3},
					},
				},
				klotski.Piece{
					Label:  "h",
					Width:  1,
					Height: 1,
					Blocks: []klotski.Block{
						klotski.Block{X: 2, Y: 3},
					},
				},
				klotski.Piece{
					Label:  "i",
					Width:  1,
					Height: 1,
					Blocks: []klotski.Block{
						klotski.Block{X: 0, Y: 4},
					},
				},
				klotski.Piece{
					Label:  "j",
					Width:  1,
					Height: 1,
					Blocks: []klotski.Block{
						klotski.Block{X: 3, Y: 4},
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
