package physics

import (
	"github.com/faiface/pixel"
)

type Speed struct {
	X float64
	Y float64
}

func Contains(object1 pixel.Rect, object2 pixel.Rect) bool {
	if object1.Max.X >= object2.Max.X &&
		object1.Max.Y >= object2.Max.Y &&
		object1.Min.X <= object2.Min.X &&
		object1.Min.Y <= object2.Min.Y {
		return true
	}

	return false
}

func Intersects(object1 pixel.Rect, object2 pixel.Rect) bool {
	intersect := object1.Intersect(object2)

	if intersect.H() == 0 && intersect.W() == 0 {
		return false
	}

	return true
}

func CalculateDiff(object pixel.Rect, immobile pixel.Rect) (repositionVec pixel.Vec) {
	if object.Min.X < immobile.Min.X { // object is to the left
		repositionVec.X = immobile.Min.X - object.Max.X
	} else if object.Max.X > immobile.Max.X { // object is to the right
		repositionVec.X = immobile.Max.X - object.Min.X
	}

	if object.Min.Y < immobile.Min.Y { // object is under
		repositionVec.Y = immobile.Min.Y - object.Max.Y
	} else if object.Max.Y > immobile.Max.Y { // object is above
		repositionVec.Y = immobile.Max.Y - object.Min.Y
	}

	return repositionVec
}

func CalculateBoundsRestriction(object pixel.Rect, bounds pixel.Rect) (speedMult Speed, repositionVec pixel.Vec) {
	speedMult = Speed{
		X: 1,
		Y: 1,
	}
	repositionVec = pixel.Vec{}

	if object.Min.X < bounds.Min.X {
		repositionVec.X = bounds.Min.X - object.Min.X
		speedMult.X = -1
	} else if object.Max.X > bounds.Max.X {
		repositionVec.X = bounds.Max.X - object.Max.X
		speedMult.X = -1
	}

	if object.Min.Y < bounds.Min.Y {
		repositionVec.Y = bounds.Min.Y - object.Min.Y
		speedMult.Y = -1
	} else if object.Max.Y > bounds.Max.Y {
		repositionVec.Y = bounds.Max.Y - object.Max.Y
		speedMult.Y = -1
	}

	return
}
