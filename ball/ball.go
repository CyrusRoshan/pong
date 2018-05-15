package ball

import (
	"github.com/CyrusRoshan/pong/physics"
	"github.com/CyrusRoshan/pong/utils"
	"github.com/faiface/pixel"
)

const (
	SPRITE = "sprites/ball.png"

	// values under 1 decelerate, values above 1, accelerate
	WALL_ACCEL     = 0.8
	CONSTANT_ACCEL = 0.95
)

type Ball struct {
	sheet  pixel.Picture
	frame  pixel.Rect
	sprite *pixel.Sprite

	Rect  pixel.Rect
	Speed physics.Speed
}

func NewBall(location pixel.Vec) *Ball {
	pic, err := utils.LoadPicture(SPRITE)
	if err != nil {
		panic(err)
	}

	b := Ball{
		sheet:  pic,
		frame:  pixel.Rect{},
		sprite: pixel.NewSprite(pic, pic.Bounds()),

		Rect: pixel.Rect{},
		Speed: physics.Speed{
			X: float64(0),
			Y: float64(0),
		},
	}

	b.Rect = b.sprite.Frame().
		Moved(location).
		Moved(pixel.Vec{X: 0, Y: -b.sprite.Frame().Size().Y / 2})

	return &b
}

func (b *Ball) Draw(t pixel.Target) {
	b.sprite.Draw(t, pixel.IM.Moved(b.Rect.Center()))
}

func (b *Ball) CalculateMovement(dt float64) {
	b.Speed.X *= CONSTANT_ACCEL
	b.Speed.Y *= CONSTANT_ACCEL

	b.Rect = b.Rect.Moved(pixel.Vec{
		X: b.Speed.X,
		Y: b.Speed.Y,
	})
}
