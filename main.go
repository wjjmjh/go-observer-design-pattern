package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	width    = 10
	height   = 5
	player   = 'P'
	obstacle = 'X'
	exit     = 'E'
)

type Observer interface {
	Update(*Game)
}

type Broadcaster1 struct{}

func (p *Broadcaster1) Update(g *Game) {
	fmt.Printf("omg!!! player moved!!! new position is (%d, %d)!\n", g.playerX, g.playerY)
}

type Broadcaster2 struct{}

func (p *Broadcaster2) Update(g *Game) {
	fmt.Printf("Alert! Alert! The player landed at coordinates (%d, %d)!\n", g.playerX, g.playerY)
}

type Game struct {
	grid      []string
	playerX   int
	playerY   int
	exitX     int
	exitY     int
	observers []Observer
}

func (g *Game) Initialize() {
	g.grid = make([]string, height)

	// create an empty row
	emptyRow := ""
	for x := 0; x < width; x++ {
		emptyRow += " "
	}

	// fill the grid with empty rows
	for y := 0; y < height; y++ {
		g.grid[y] = emptyRow
	}

	// initial player position
	g.playerX = 0
	g.playerY = 0
	g.grid[g.playerY] = replaceAtIndex(g.grid[g.playerY], player, g.playerX)

	// exit position
	g.exitX = width - 1
	g.exitY = height - 1
	g.grid[g.exitY] = replaceAtIndex(g.grid[g.exitY], exit, g.exitX)

	// add some obstacles
	g.grid[1] = replaceAtIndex(g.grid[1], obstacle, 2)
	g.grid[2] = replaceAtIndex(g.grid[2], obstacle, 4)
	g.grid[3] = replaceAtIndex(g.grid[3], obstacle, 1)

	g.observers = append(g.observers, &Broadcaster1{})
	g.observers = append(g.observers, &Broadcaster2{})
}

func (g *Game) NotifyObservers() {
	for _, observer := range g.observers {
		observer.Update(g)
	}
}

func (g *Game) drawGame() {
	for _, row := range g.grid {
		fmt.Println(row)
	}
}

func (g *Game) handleInput() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter movement (w/a/s/d): ")
	input, _ := reader.ReadString('\n')
	input = input[:1] // Extract the first character
	g.movePlayer(input)
}

func (g *Game) movePlayer(input string) {
	newX, newY := g.playerX, g.playerY

	switch input {
	case "w":
		newY--
	case "s":
		newY++
	case "a":
		newX--
	case "d":
		newX++
	}

	// check if new position is valid
	if newX >= 0 && newX < width && newY >= 0 && newY < height && g.grid[newY][newX] != obstacle {
		// clear previous player position
		g.grid[g.playerY] = replaceAtIndex(g.grid[g.playerY], ' ', g.playerX)

		// update player position
		g.playerX, g.playerY = newX, newY
		g.grid[g.playerY] = replaceAtIndex(g.grid[g.playerY], player, g.playerX)

		// notify broadcasters
		g.NotifyObservers()

		// check if the player reached the exit
		if g.playerX == g.exitX && g.playerY == g.exitY {
			fmt.Println("Congratulations! You reached the exit.")
			os.Exit(0)
		}
	}
}

func replaceAtIndex(s string, r rune, index int) string {
	if index < 0 || index >= len(s) {
		return s
	}
	return s[:index] + string(r) + s[index+1:]
}

func main() {
	game := Game{}
	game.Initialize()

	for {
		game.drawGame()
		game.handleInput()
	}
}
