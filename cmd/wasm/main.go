package main

import (
	"math/rand"
	"syscall/js"
	"time"
)

type Point struct {
	X, Y int
}

type Snake struct {
	Body      []Point
	Direction Point
}

type Apple struct {
	Pos Point
}

type GameState struct {
	Snake    Snake
	Apple    Apple
	GameOver bool
	Score    int
}

var (
	state         GameState
	ctx           js.Value
	scoreEl       js.Value
	width, height int
)

const (
	CellSize = 20
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	// Init Canvas stuff
	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "canvas")
	width = canvasEl.Get("clientWidth").Int() / CellSize
	height = canvasEl.Get("clientHeight").Int() / CellSize
	canvasEl.Set("width", width*CellSize)
	canvasEl.Set("height", height*CellSize)
	ctx = canvasEl.Call("getContext", "2d")
	scoreEl = doc.Call("getElementById", "score")

	keyDownEvt := js.FuncOf(keyDownEvent)
	defer keyDownEvt.Release()
	doc.Call("addEventListener", "keydown", keyDownEvt)

	reset()
	gameLoop()
}

func keyDownEvent(this js.Value, args []js.Value) interface{} {
	e := args[0]
	keyCode := e.Get("keyCode").Int()

	switch keyCode {
	case 37: // left arrow
		state.Snake.Direction = Point{-1, 0}
	case 38: // up arrow
		state.Snake.Direction = Point{0, -1}
	case 39: // right arrow
		state.Snake.Direction = Point{1, 0}
	case 40: // down arrow
		state.Snake.Direction = Point{0, 1}
	}

	return nil
}

func reset() {
	state = GameState{
		Snake: Snake{
			Body:      []Point{{X: width / 2, Y: height / 2}},
			Direction: Point{X: 1, Y: 0},
		},
		Apple: Apple{
			Pos: Point{X: rand.Intn(width), Y: rand.Intn(height)},
		},
		GameOver: false,
	}
}

func gameLoop() {
	for {
		if state.GameOver {
			js.Global().Call("handleGameOver", state.Score)
			break
		}
		update()
		draw()
		time.Sleep(getDelay((state.Score)))
	}
}

func getDelay(score int) time.Duration {
	delay := time.Second / time.Duration(10+score)
	return delay
}

func update() {
	head := state.Snake.Body[0]
	next := Point{
		X: (head.X + state.Snake.Direction.X + width) % width,
		Y: (head.Y + state.Snake.Direction.Y + height) % height,
	}

	if collidesWithSnake(next) {
		state.GameOver = true
		return
	}

	state.Snake.Body = append([]Point{next}, state.Snake.Body...)

	if collidesWithApple(next) {
		state.Apple.Pos = Point{X: rand.Intn(width), Y: rand.Intn(height)}
		state.Score++
		scoreEl.Set("innerHTML", state.Score)
	} else {
		state.Snake.Body = state.Snake.Body[:len(state.Snake.Body)-1]
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
	ctx.Call("clearRect", 0, 0, width*CellSize, height*CellSize)

	ctx.Set("fillStyle", "red")
	drawPoint(state.Apple.Pos)

	ctx.Set("fillStyle", "green")
	for _, p := range state.Snake.Body {
		drawPoint(p)
	}
}

func drawPoint(p Point) {
	ctx.Call("fillRect", p.X*CellSize, p.Y*CellSize, CellSize, CellSize)
}
