package gameholder

import (
	"time"

	"github.com/CyrusRoshan/pong/physics"

	"github.com/CyrusRoshan/pong/ball"
	"github.com/CyrusRoshan/pong/player"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type GameHolder struct {
	players    []*player.Player
	ball       *ball.Ball
	lastRender time.Time
}

func MakeGame(canvas *pixelgl.Canvas, totalPlayers int) *GameHolder {
	gh := GameHolder{}
	bounds := canvas.Bounds()

	// make ball
	gh.ball = ball.NewBall(bounds.Center())

	// make players
	gh.players = make([]*player.Player, totalPlayers)
	for i := 0; i < totalPlayers; i++ {
		xPos := bounds.Center().X
		if i%2 == 1 {
			xPos += bounds.W() / 4
		} else {
			xPos -= bounds.W() / 4
		}

		location := pixel.Vec{
			X: xPos,
			Y: bounds.Center().Y,
		}

		gh.players[i] = player.NewPlayer((i%2 == 1), i, location)
	}

	return &gh
}

func (gh *GameHolder) GetChanges(win *pixelgl.Window) {
	dt := gh.timeSinceLastRender()

	// calc ball movement
	gh.ball.CalculateMovement(dt)

	// calc player movement
	for i := 0; i < len(gh.players); i++ {
		p := gh.players[i]

		var keypairs []player.KeyPair

		// only 1 to 2-player key support right now
		if i == 0 {
			keypairs = []player.KeyPair{
				player.KeyPair{
					Key:        pixelgl.KeyW,
					MirrorKey:  pixelgl.KeyS,
					IsVertical: true,
					SpeedVar:   &p.Speed.Y,
					AccelVar:   player.MOVE_ACCEL,
				},
				player.KeyPair{
					Key:        pixelgl.KeyD,
					MirrorKey:  pixelgl.KeyA,
					IsVertical: false,
					SpeedVar:   &p.Speed.X,
					AccelVar:   player.MOVE_ACCEL,
				},
			}
		} else {
			keypairs = []player.KeyPair{
				player.KeyPair{
					Key:        pixelgl.KeyUp,
					MirrorKey:  pixelgl.KeyDown,
					IsVertical: true,
					SpeedVar:   &p.Speed.Y,
					AccelVar:   player.MOVE_ACCEL,
				},
				player.KeyPair{
					Key:        pixelgl.KeyRight,
					MirrorKey:  pixelgl.KeyLeft,
					IsVertical: false,
					SpeedVar:   &p.Speed.X,
					AccelVar:   player.MOVE_ACCEL,
				},
			}
		}

		p.GetInput(win, keypairs, dt)
	}
}

func (gh *GameHolder) RestrictBoundsTo(bounds pixel.Rect) {
	speedMult, movedVec := physics.RestrictBoundsTo(gh.ball.Rect, bounds, ball.WALL_ACCEL)
	gh.ball.Speed.X *= speedMult.X
	gh.ball.Speed.Y *= speedMult.Y

	gh.ball.Rect = gh.ball.Rect.Moved(movedVec)

	for i := 0; i < len(gh.players); i++ {
		speedMult, movedVec := physics.RestrictBoundsTo(gh.players[i].Rect, bounds, player.WALL_ACCEL)
		gh.players[i].Speed.X *= speedMult.X
		gh.players[i].Speed.Y *= speedMult.Y

		gh.players[i].Rect = gh.players[i].Rect.Moved(movedVec)
	}
}

func (gh *GameHolder) Draw(t pixel.Target) {
	gh.ball.Draw(t)

	for i := 0; i < len(gh.players); i++ {
		gh.players[i].Draw(t)
	}
}

func (gh *GameHolder) timeSinceLastRender() float64 {
	diff := time.Since(gh.lastRender).Seconds()
	gh.lastRender = time.Now()

	return diff
}
