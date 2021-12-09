package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"math/rand"
	"strconv"
)

type Coordinate struct {
	y int
	x int
}

// board cell values
const (
	Empty int = iota
	Body 
	Head
	Food
	DeadHead
)

var default_height, default_width = 10, 10 //default board dimensions

type GameState struct {
	Board [][]int
	BoardDimension [2]int
	Round int
	Score int
	SnakeLength int
	SnakeHead Coordinate
	SnakeBody []Coordinate
	Heading string
	Food Coordinate
}


func main() {

	state := InitializeState()
	DisplayState(state)
	move := RequireMove()

	for {
		if ok := UpdateState(&state, move); !ok {
			DisplayState(state)
			fmt.Printf("GAME OVER! Final Score: %d\n", state.Score)
			return 
		}

		DisplayState(state)
		move = RequireMove()
	}
}

// user inputs move direction W A S D 
// or inputs nothing to keep going in current direction
// returns "W" "A" "S" "D" or ""
func RequireMove() string {

	fmt.Printf("Input direction and hit enter: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	user_choice := strings.ToUpper(input.Text())

	if user_choice == "W" || user_choice == "S" || user_choice == "A" || user_choice == "D" {
		return user_choice
	}

	return ""
}


// takes dimensions as input and returns initial game state
func InitializeState() GameState {

	var height, width int

	if len(os.Args) < 3 {
		fmt.Println("using default dimensions...")
		height, width = default_height, default_width
	} else {
		var err1, err2 error
		height, err1 = strconv.Atoi(os.Args[1])
		width, err2 = strconv.Atoi(os.Args[2])

		if err1 != nil || err2 != nil {
			fmt.Println("using default dimensions...")
			height, width = default_height, default_width
		} else {
			fmt.Printf("using %d by %d board dimensions\n", height, width)
		}
	}

	state := GameState {
		BoardDimension: [2]int{height, width},
		Round: 0,
		Score: 0,
		SnakeLength: 2,
		SnakeHead: Coordinate{height/2, width/2},
		SnakeBody: []Coordinate{Coordinate{height/2, width/2-1}},
		Heading: "D",
	}

	state.Board = make([][]int, height)
	for i := 0; i < height; i++ {
		state.Board[i] = make([]int, width)
	}

	// get random food location
	ok := GetNewFoodLoc(&state)
	if !ok {
		fmt.Println("could not get food loc while initializing")
	}

	// populate board 

	state.Board[state.SnakeHead.y][state.SnakeHead.x] = Head 
	state.Board[state.SnakeBody[0].y][state.SnakeBody[0].x] = Body

	return state
}

// takes current state and move as input
// returns new game state and bool ok which is false for game over move
func UpdateState(state *GameState, move string) bool {
	if move == "" {
		move = state.Heading
	}

	prevHead := state.SnakeHead

	switch move {
	case "W":
		if state.Heading == "S" {
			return true
		}
		state.SnakeHead.y--
	case "S":
		if state.Heading == "W" {
			return true
		}
		state.SnakeHead.y++
	case "D":
		if state.Heading == "A" {
			return true
		}
		state.SnakeHead.x++
	case "A":
		if state.Heading == "D" {
			return true
		}
		state.SnakeHead.x--
	}

	// check if game over
	if state.SnakeHead.y >= state.BoardDimension[0] || state.SnakeHead.y < 0 || state.SnakeHead.x >= state.BoardDimension[1] || state.SnakeHead.x < 0 {
		state.Board[prevHead.y][prevHead.x] = DeadHead
		return false
	}
	if state.Board[state.SnakeHead.y][state.SnakeHead.x] == Body {
		state.Board[state.SnakeHead.y][state.SnakeHead.x] = DeadHead
		return false
	}	
	
	state.Heading = move

	// move body
	state.SnakeBody = append(state.SnakeBody, prevHead)

	// check for food and move head
	if state.Board[state.SnakeHead.y][state.SnakeHead.x] == Food {
		state.Score++
		ok := GetNewFoodLoc(state)
		if !ok {
			fmt.Println("could not get food loc")
		}
	} else {
		state.Board[state.SnakeBody[0].y][state.SnakeBody[0].x] = Empty
		state.SnakeBody = state.SnakeBody[1:len(state.SnakeBody)]
	}

	state.Board[state.SnakeHead.y][state.SnakeHead.x] = Head
	state.Board[prevHead.y][prevHead.x] = Body

	state.Round++

	return true
	
}

// takes state as input and prints board to console
func DisplayState(state GameState) {
	top_border := " "
	for i := 0; i < state.BoardDimension[1]; i++ {
		top_border += "_"
	}
	fmt.Printf("%s \n", top_border)
	for i, row := range state.Board {
		fmt.Printf("|")
		for _, col := range row {
			switch col {
			case Empty:
				fmt.Printf(" ")
			case Body:
				fmt.Printf("#")
			case Food:
				fmt.Printf("o")
			case Head:
				fmt.Printf("@")
			case DeadHead:
				fmt.Printf("X")
			}
		}
		fmt.Printf("|")
		if i == 0 {fmt.Printf("\tRound: %d", state.Round)}
		if i == 1 {fmt.Printf("\tScore: %d", state.Score)}
		if i == 2 {fmt.Printf("\tControls: W: up |A: left |S: down |D: right |other/nothing: move on same heading")}
		fmt.Printf("\n")
	}
	bottom_border := " "
	for i := 0; i < state.BoardDimension[1]; i++ {
		bottom_border += "*"
	}
	fmt.Printf("%s \n", bottom_border)
}

// takes as input game state 
// returns a random food location
func GetNewFoodLoc(state *GameState) bool {
	h, w := state.BoardDimension[0], state.BoardDimension[1]

	for count := 0; count < 10*h*w; count++{
		y, x := rand.Intn(h), rand.Intn(w)
		
		if state.Board[y][x] == Empty {
			state.Food = Coordinate{y, x}
			state.Board[y][x] = Food
			return true
		}
	}
	return false
}