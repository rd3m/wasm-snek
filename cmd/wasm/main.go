package main

import (
	"syscall/js"
)

func main() {
	// Init Canvas stuff
	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "canvas")
	width := doc.Get("body").Get("clientWidth").Float()
	height := doc.Get("body").Get("clientHeight").Float()
	canvasEl.Set("width", width)
	canvasEl.Set("height", height)
	canvasEl.Set("style", "border: 1px solid black;")
	// ctx := canvasEl.Call("getContext", "2d")
}
