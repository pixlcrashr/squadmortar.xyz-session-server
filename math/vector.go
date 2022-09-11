package math

type Vector3 struct {
	X float32
	Y float32
	Z float32
}

func (v Vector3) Add(value Vector3) Vector3 {
	return Vector3{
		v.X + value.X,
		v.Y + value.Y,
		v.Z + value.Z,
	}
}
