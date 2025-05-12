package gameEntities

type InterestPoint uint8

const (
	Food InterestPoint = iota
	Structure
	OtherCreature
)

type Point struct {
	X, Y  float32
	PType InterestPoint
}

func (p Point) Coord() (x, y float32) {
	return p.X, p.Y
}

func (p Point) Type() InterestPoint {
	return p.PType
}
