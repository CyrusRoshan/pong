package player

import (
	"math"

	"github.com/CyrusRoshan/pong/physics"
	"github.com/CyrusRoshan/pong/utils"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	SPRITE = "sprites/paddle.png"

	MOVE_ACCEL          = 200
	BOOST_REVERSE_ACCEL = 3

	CONSTANT_ACCEL = 0.5
	STOP_ACCEL     = 5
	WALL_ACCEL     = 0.3

	MAX_SPEED = float64(80)
)

type Player struct {
	playerNum int
	direction float64

	sheet  pixel.Picture
	frame  pixel.Rect
	sprite *pixel.Sprite

	Rect  pixel.Rect
	Speed physics.Speed
}

func NewPlayer(isTeamTwo bool, playerNumber int, location pixel.Vec) *Player {
	pic, err := utils.LoadPicture(SPRITE)
	if err != nil {
		panic(err)
	}

	direction := 1
	if isTeamTwo {
		direction = -1
	}

	p := Player{
		playerNum: playerNumber,
		direction: float64(direction),

		sheet:  pic,
		frame:  pixel.Rect{},
		sprite: pixel.NewSprite(pic, pic.Bounds()),

		Rect: pixel.Rect{},
		Speed: physics.Speed{
			X: float64(0),
			Y: float64(0),
		},
	}

	p.Rect = p.sprite.Frame().
		Moved(location).
		Moved(pixel.Vec{X: 0, Y: -p.sprite.Frame().Size().Y / 2})

	return &p
}

func (p *Player) Draw(t pixel.Target) {
	p.sprite.Draw(t,
		pixel.IM.
			ScaledXY(
				pixel.ZV,
				pixel.Vec{X: p.direction, Y: 1}).
			Moved(p.Rect.Center()))
}

type KeyPair struct {
	Key        pixelgl.Button
	MirrorKey  pixelgl.Button
	IsVertical bool
	SpeedVar   *float64
	AccelVar   int
}

func (p *Player) GetInput(win *pixelgl.Window, Keypairs []KeyPair, dt float64) {
	movedVector := pixel.Vec{}

	for i := 0; i < len(Keypairs); i++ {
		pair := Keypairs[i]

		ds := dt * float64(pair.AccelVar)
		keyPressed := win.Pressed(pair.Key)
		mirrorkeyPressed := win.Pressed(pair.MirrorKey)

		accelerate := keyPressed || mirrorkeyPressed
		if keyPressed && mirrorkeyPressed {
			accelerate = false
		}

		if accelerate {
			if keyPressed {
				if *pair.SpeedVar < 0 {
					ds *= BOOST_REVERSE_ACCEL
				}

				*pair.SpeedVar += ds
			} else { // if mirrorkeyPressed
				if *pair.SpeedVar > 0 {
					ds *= BOOST_REVERSE_ACCEL
				}

				*pair.SpeedVar -= ds
			}

			if *pair.SpeedVar > MAX_SPEED {
				*pair.SpeedVar = MAX_SPEED
			} else if *pair.SpeedVar < -MAX_SPEED {
				*pair.SpeedVar = -MAX_SPEED
			}
		} else {
			if *pair.SpeedVar > 0 {
				*pair.SpeedVar -= ds * CONSTANT_ACCEL
			} else if *pair.SpeedVar < 0 {
				*pair.SpeedVar += ds * CONSTANT_ACCEL
			}

			if math.Abs(*pair.SpeedVar) < STOP_ACCEL {
				*pair.SpeedVar = 0
			}
		}

		ddist := *pair.SpeedVar * dt

		if pair.IsVertical {
			movedVector.Y += ddist
		} else {
			movedVector.X += ddist
		}
	}

	p.Rect = p.Rect.Moved(movedVector)
}
