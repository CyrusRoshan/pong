package gameholder

import (
	"math"
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

func (gh *GameHolder) CalculateCollisions() {
	for i := 0; i < len(gh.players); i++ {
		if physics.Intersects(gh.players[i].Rect, gh.ball.Rect) {
			diff := physics.CalculateDiff(gh.ball.Rect, gh.players[i].Rect)

			if math.Abs(gh.players[i].Speed.X) > math.Abs(gh.players[i].Speed.Y) { // only move either X or Y
				diff.Y = 0
			} else {
				diff.X = 0
			}
			gh.ball.Rect = gh.ball.Rect.Moved(diff)

			multVec := physics.CalculateCollisionVector(gh.players[i].Rect, gh.ball.Rect)

			// gh.ball.Speed.X += multVec.X
			// if multVec.X == -1 {
			// 	gh.ball.Speed.X += math.Copysign(ball.PLAYER_HIT_SPEED_BOOST, gh.ball.Speed.X)
			// }

			// gh.ball.Speed.Y *= multVec.Y
			// if multVec.Y == -1 {
			// 	gh.ball.Speed.Y += math.Copysign(ball.PLAYER_HIT_SPEED_BOOST, gh.ball.Speed.Y)
			// }
		}
	}
}

func (gh *GameHolder) RestrictBoundsTo(bounds pixel.Rect) {
	if !physics.Contains(bounds, gh.ball.Rect) {
		speedMult, movedVec := physics.CalculateBoundsRestriction(gh.ball.Rect, bounds)
		gh.ball.Speed.X *= speedMult.X * ball.WALL_ACCEL
		gh.ball.Speed.Y *= speedMult.Y * ball.WALL_ACCEL

		gh.ball.Rect = gh.ball.Rect.Moved(movedVec)
	}

	for i := 0; i < len(gh.players); i++ {
		if !physics.Contains(bounds, gh.players[i].Rect) {
			speedMult, movedVec := physics.CalculateBoundsRestriction(gh.players[i].Rect, bounds)
			gh.players[i].Speed.X *= speedMult.X * player.WALL_ACCEL
			gh.players[i].Speed.Y *= speedMult.Y * player.WALL_ACCEL

			gh.players[i].Rect = gh.players[i].Rect.Moved(movedVec)
		}
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
