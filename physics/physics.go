package physics

import "github.com/faiface/pixel"

type Speed struct {
	X float64
	Y float64
}

func RestrictBoundsTo(object pixel.Rect, bounds pixel.Rect, collisionAccel float64) (speedMult Speed, movedVec pixel.Vec) {
	speedMult = Speed{
		X: 1,
		Y: 1,
	}
	movedVec = pixel.Vec{}

	if object.Min.X < bounds.Min.X {
		movedVec.X = bounds.Min.X - object.Min.X
		speedMult.X = -collisionAccel
	} else if object.Max.X > bounds.Max.X {
		movedVec.X = bounds.Max.X - object.Max.X
		speedMult.X = -collisionAccel
	}

	if object.Min.Y < bounds.Min.Y {
		movedVec.Y = bounds.Min.Y - object.Min.Y
		speedMult.Y = -collisionAccel
	} else if object.Max.Y > bounds.Max.Y {
		movedVec.Y = bounds.Max.Y - object.Max.Y
		speedMult.Y = -collisionAccel
	}

	return
}
