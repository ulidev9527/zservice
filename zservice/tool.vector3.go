package zservice

import "math"

type Vector3 struct {
	X, Y, Z float64
}

func (v Vector3) Distance(v2 Vector3) float64 {
	dx := v.X - v2.X
	dy := v.Y - v2.Y
	dz := v.Z - v2.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}
