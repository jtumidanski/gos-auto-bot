package coordinate

type Model struct {
	x     int
	y     int
	scale int
}

func New(x int, y int) Model {
	return Model{
		x:     x,
		y:     y,
		scale: 1,
	}
}

func NewScaled(x int, y int) Model {
	return Model{
		x:     x,
		y:     y,
		scale: 1,
	}
}

func (m *Model) X() int {
	return m.x * m.scale
}

func (m *Model) Y() int {
	return m.y * m.scale
}
