package main

import (
	"math"

	"github.com/CyrusRoshan/pongg/player"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"golang.org/x/image/colornames"
)

func main() {
	pixelgl.Run(run)
}

func run() {
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

	const SCALE = 15
	canvas := pixelgl.NewCanvas(pixel.R(-width/SCALE, -height/SCALE, width/SCALE, height/SCALE))
	cam := pixel.IM
	canvas.SetMatrix(cam)

	// create players
	players := player.MakePlayers(canvas, 2)

	for !win.Closed() {
		win.SetMatrix(pixel.IM.Scaled(pixel.ZV,
			math.Min(
				win.Bounds().W()/canvas.Bounds().W(),
				win.Bounds().H()/canvas.Bounds().H(),
			),
		).Moved(win.Bounds().Center()))

		win.Clear(colornames.Red)
		canvas.Clear(colornames.Black)

		// player get input
		// player draw
		players.GetInput(win)
		players.Draw(canvas)

		canvas.Draw(win, pixel.IM.Moved(canvas.Bounds().Center()))

		win.Update()
	}
}
