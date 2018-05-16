package main

import (
	"math"

	"github.com/CyrusRoshan/pong/gameholder"
	"github.com/CyrusRoshan/pong/utils"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"golang.org/x/image/colornames"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	// set up window
	width, height := pixelgl.PrimaryMonitor().Size()

	cfg := pixelgl.WindowConfig{
		Title:     "Pongg",
		Bounds:    pixel.R(0, 0, width/2, height/2),
		VSync:     true,
		Resizable: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.Clear(colornames.Red)

	// make canvas (drawing to canvas, then scaling to window)
	const SCALE = 15
	var (
		canvas = pixelgl.NewCanvas(pixel.R(-width/SCALE, -height/SCALE, width/SCALE, height/SCALE))
		cam    = pixel.IM
	)
	canvas.SetMatrix(cam)

	// create players
	game := gameholder.MakeGame(canvas, 2)

	// game loop
	for !win.Closed() {
		// make canvas fit window
		win.SetMatrix(pixel.IM.Scaled(pixel.ZV,
			math.Min(
				win.Bounds().W()/canvas.Bounds().W(),
				win.Bounds().H()/canvas.Bounds().H(),
			),
		).Moved(win.Bounds().Center()))

		// get input, change model
		game.GetChanges(win)
		game.RestrictBoundsTo(canvas.Bounds())
		game.CalculateCollisions()

		// clear window before redrawing
		win.Clear(colornames.Red)
		canvas.Clear(colornames.Black)

		// draw to canvas -> window -> update
		game.Draw(canvas)
		canvas.Draw(win, pixel.IM.Moved(canvas.Bounds().Center()))
		win.Update()

		utils.DisplayFPS(win, cfg.Title)
	}
}
