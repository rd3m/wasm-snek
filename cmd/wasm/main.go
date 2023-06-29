package main

import (
	"math/rand"
	"syscall/js"
	"time"
)

type Point struct {
	X, Y float64
}

type Snake struct {
	Body    []Point
	Dir     Point // The direction the snake is moving
	Growing bool  // Whether the snake is growing
}

type Apple struct {
	Pos Point
}

type GameState struct {
	Snake    Snake
	Apple    Apple
	GameOver bool
}

var state GameState
var ctx js.Value

var width, height float64

const (
	CellSize = 20
	Speed    = time.Second / 10
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	// Init Canvas stuff
	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "canvas")
	width = canvasEl.Get("clientWidth").Float()
	height = canvasEl.Get("clientHeight").Float()
	canvasEl.Set("width", width)
	canvasEl.Set("height", height)
	ctx = canvasEl.Call("getContext", "2d")

	reset()
	gameLoop()
}

func reset() {
	state = GameState{
		Snake: Snake{
			Body: []Point{{X: width / 2, Y: height / 2}},
			Dir:  Point{X: CellSize, Y: 0},
		},
		Apple: Apple{
			Pos: Point{X: rand.Float64() * width, Y: rand.Float64() * height},
		},
		GameOver: false,
	}
}

func gameLoop() {
	ticker := time.NewTicker(Speed)
	defer ticker.Stop()

	for range ticker.C {
		update()
		draw()
	}
}

func update() {
	if state.GameOver {
		return
	}

	head := state.Snake.Body[0]
	next := Point{
		X: head.X + state.Snake.Dir.X,
		Y: head.Y + state.Snake.Dir.Y,
	}

	if next.X < 0 || next.Y < 0 || next.X >= width || next.Y >= height {
		state.GameOver = true
		return
	}

	if collidesWithSnake(next) {
		state.GameOver = true
		return
	}

	state.Snake.Body = append([]Point{next}, state.Snake.Body...)

	if !collidesWithApple(next) {
		state.Snake.Body = state.Snake.Body[:len(state.Snake.Body)-1]
	} else {
		state.Apple.Pos = Point{X: rand.Float64() * width, Y: rand.Float64() * height}
	}
}

func collidesWithSnake(p Point) bool {
	for _, b := range state.Snake.Body {
		if p.X == b.X && p.Y == b.Y {
			return true
		}
	}

	return false
}

func collidesWithApple(p Point) bool {
	return p.X == state.Apple.Pos.X && p.Y == state.Apple.Pos.Y
}

func draw() {
	// Clear the canvas
	ctx.Call("clearRect", 0, 0, width, height)

	// Draw the apple
	ctx.Set("fillStyle", "red")
	drawPoint(state.Apple.Pos)

	// Draw the snake
	ctx.Set("fillStyle", "green")
	for _, p := range state.Snake.Body {
		drawPoint(p)
	}

	// If the game is over, draw a message
	if state.GameOver {
		ctx.Set("font", "48px serif")
		ctx.Set("fillStyle", "black")
		ctx.Call("fillText", "Game Over", width/3, height/2)
	}
}

func drawPoint(p Point) {
	ctx.Call("fillRect", p.X, p.Y, CellSize, CellSize)
}
