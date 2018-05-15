package player

import (
	"math"
	"time"

	"github.com/CyrusRoshan/pong/utils"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const PLAYER_SPRITE = "sprites/paddle.png"

type PlayerHolder struct {
	players []*player
}

func MakePlayers(canvas *pixelgl.Canvas, totalPlayers int) *PlayerHolder {
	ph := PlayerHolder{}

	ph.players = make([]*player, totalPlayers)
	bounds := canvas.Bounds()

	for i := 0; i < totalPlayers; i++ {
		ph.players[i] = newPlayer(PLAYER_SPRITE, (i%2 == 1), i)

		xPos := bounds.Center().X
		if i%2 == 1 {
			xPos += bounds.W() / 4
		} else {
			xPos -= bounds.W() / 4
		}

		ph.players[i].rect = ph.players[i].rect.Moved(pixel.Vec{
			X: xPos,
			Y: bounds.Center().Y,
		})
	}

	return &ph
}

func (ph *PlayerHolder) GetInput(win *pixelgl.Window) {
	for i := 0; i < len(ph.players); i++ {
		p := ph.players[i]

		var keypairs []keyPair

		if i == 0 {
			keypairs = []keyPair{
				keyPair{
					key:        pixelgl.KeyW,
					mirrorKey:  pixelgl.KeyS,
					isVertical: true,
					speedVar:   &p.ySpeed,
					accelVar:   PLAYER_ACCEL,
				},
				keyPair{
					key:        pixelgl.KeyD,
					mirrorKey:  pixelgl.KeyA,
					isVertical: false,
					speedVar:   &p.xSpeed,
					accelVar:   PLAYER_ACCEL,
				},
			}
		} else {
			keypairs = []keyPair{
				keyPair{
					key:        pixelgl.KeyUp,
					mirrorKey:  pixelgl.KeyDown,
					isVertical: true,
					speedVar:   &p.ySpeed,
					accelVar:   PLAYER_ACCEL,
				},
				keyPair{
					key:        pixelgl.KeyRight,
					mirrorKey:  pixelgl.KeyLeft,
					isVertical: false,
					speedVar:   &p.xSpeed,
					accelVar:   PLAYER_ACCEL,
				},
			}
		}

		p.getInput(win, keypairs)
	}
}

func (ph *PlayerHolder) RestrictBoundsTo(bounds pixel.Rect) {
	for i := 0; i < len(ph.players); i++ {
		ph.players[i].restrictBoundsTo(bounds)
	}
}

func (ph *PlayerHolder) Draw(t pixel.Target) {
	for i := 0; i < len(ph.players); i++ {
		ph.players[i].draw(t)
	}
}

const PLAYER_ACCEL = 200
const PLAYER_MAX_SPEED = float64(80)

type player struct {
	playerNum int
	sheet     pixel.Picture

	frame  pixel.Rect
	sprite *pixel.Sprite

	direction float64
	rect      pixel.Rect

	xSpeed float64
	ySpeed float64

	lastRender time.Time
}

func newPlayer(picLocation string, isTeamTwo bool, playerNumber int) *player {
	pic, err := utils.LoadPicture(picLocation)
	if err != nil {
		panic(err)
	}

	direction := 1
	if isTeamTwo {
		direction = -1
	}

	p := player{
		playerNum: playerNumber,
		sheet:     pic,

		xSpeed: float64(0),
		ySpeed: float64(0),

		rect:      pixel.Rect{},
		direction: float64(direction),

		frame:  pixel.Rect{},
		sprite: pixel.NewSprite(pic, pic.Bounds()),

		lastRender: time.Now(),
	}

	p.rect = p.sprite.Frame()

	return &p
}

func (p *player) draw(t pixel.Target) {
	p.sprite.Draw(t,
		pixel.IM.
			ScaledXY(
				pixel.ZV,
				pixel.Vec{X: p.direction, Y: 1}).
			Moved(p.rect.Center()))
}

type keyPair struct {
	key        pixelgl.Button
	mirrorKey  pixelgl.Button
	isVertical bool
	speedVar   *float64
	accelVar   int
}

func (p *player) getInput(win *pixelgl.Window, keypairs []keyPair) {
	dt := p.timeSinceLastRender()

	movedVector := pixel.Vec{}

	for i := 0; i < len(keypairs); i++ {
		pair := keypairs[i]

		ds := dt * float64(pair.accelVar)
		keyPressed := win.Pressed(pair.key)
		mirrorkeyPressed := win.Pressed(pair.mirrorKey)

		accelerate := keyPressed || mirrorkeyPressed
		if keyPressed && mirrorkeyPressed {
			accelerate = false
		}

		if accelerate {
			if keyPressed {
				if *pair.speedVar < 0 {
					ds *= 4
				}

				*pair.speedVar += ds
			} else { // if mirrorkeyPressed
				if *pair.speedVar > 0 {
					ds *= 4
				}

				*pair.speedVar -= ds
			}

			if *pair.speedVar > PLAYER_MAX_SPEED {
				*pair.speedVar = PLAYER_MAX_SPEED
			} else if *pair.speedVar < -PLAYER_MAX_SPEED {
				*pair.speedVar = -PLAYER_MAX_SPEED
			}
		} else {
			if *pair.speedVar > 0 {
				*pair.speedVar -= (ds / 2)
			} else if *pair.speedVar < 0 {
				*pair.speedVar += (ds / 2)
			}

			if math.Abs(*pair.speedVar) < 5 {
				*pair.speedVar = 0
			}
		}

		ddist := *pair.speedVar * dt

		if pair.isVertical {
			movedVector.Y += ddist
		} else {
			movedVector.X += ddist
		}
	}

	p.rect = p.rect.Moved(movedVector)
}

func (p *player) restrictBoundsTo(bounds pixel.Rect) {
	movedVec := pixel.Vec{}

	if p.rect.Min.X < bounds.Min.X {
		movedVec.X = bounds.Min.X - p.rect.Min.X
		p.xSpeed = 0
	} else if p.rect.Max.X > bounds.Max.X {
		movedVec.X = bounds.Max.X - p.rect.Max.X
		p.xSpeed = 0
	}

	if p.rect.Min.Y < bounds.Min.Y {
		movedVec.Y = bounds.Min.Y - p.rect.Min.Y
		p.ySpeed = 0
	} else if p.rect.Max.Y > bounds.Max.Y {
		movedVec.Y = bounds.Max.Y - p.rect.Max.Y
		p.ySpeed = 0
	}

	p.rect = p.rect.Moved(movedVec)
}

func (p *player) timeSinceLastRender() float64 {
	diff := time.Since(p.lastRender).Seconds()
	p.lastRender = time.Now()

	return diff
}
