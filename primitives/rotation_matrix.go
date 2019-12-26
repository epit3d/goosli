package primitives

import "math"

type RotationMatrix struct {
	X Vector
	Y Vector
	Z Vector
}

func RotationAroundX(angle float64) RotationMatrix {
	res := RotationMatrix{}
	alpha := ToRadians(angle)
	// transposed matrix to rotate around X
	res.X = V(1, 0, 0)
	res.Y = V(0, math.Cos(alpha), math.Sin(alpha))
	res.Z = V(0, -math.Sin(alpha), math.Cos(alpha))
	return res
}

func RotationAroundZ(angle float64) RotationMatrix {
	res := RotationMatrix{}
	alpha := ToRadians(angle)
	// transposed matrix to rotate around X
	res.X = V(math.Cos(alpha), math.Sin(alpha), 0)
	res.Y = V(-math.Sin(alpha), math.Cos(alpha), 0)
	res.Z = V(0, 0, 1)
	return res
}
