package main

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type player struct {
	sheet pixel.Picture

	frame  pixel.Rect
	sprite *pixel.Sprite

	rect pixel.Rect

	xSpeed float64
	ySpeed float64

	lastRender time.Time
}

func NewPlayer(picLocation string) *player {
	pic, err := loadPicture(picLocation)
	if err != nil {
		panic(err)
	}

	return &player{
		sheet: pic,

		xSpeed: float64(20),
		ySpeed: float64(20),

		frame:  pixel.Rect{},
		sprite: pixel.NewSprite(pic, pic.Bounds()),

		lastRender: time.Now(),
	}
}

func (p *player) Draw(t pixel.Target) {
	p.sprite.Draw(t, pixel.IM.Moved(p.rect.Center()))
}

func (p *player) GetInput(win *pixelgl.Window) {
	dt := p.timeSinceLastRender()

	movedVector := pixel.Vec{}
	if win.Pressed(pixelgl.KeyW) {
		movedVector.Y += p.ySpeed * dt
	}
	if win.Pressed(pixelgl.KeyS) {
		movedVector.Y -= p.ySpeed * dt
	}
	if win.Pressed(pixelgl.KeyD) {
		movedVector.X += p.xSpeed * dt
	}
	if win.Pressed(pixelgl.KeyA) {
		movedVector.X -= p.xSpeed * dt
	}

	p.rect = p.rect.Moved(movedVector)
}

func (p *player) timeSinceLastRender() float64 {
	diff := time.Since(p.lastRender).Seconds()
	p.lastRender = time.Now()

	return diff
}
