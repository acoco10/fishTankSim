package gameEntities

type InterestPoint uint8

const (
	Food InterestPoint = iota
	Structure
	OtherCreature
	DrawPoint
)

type Point struct {
	X, Y  float32
	PType InterestPoint
}

func (p *Point) Coord() (float32, float32) {
	return p.X, p.Y
}

func (p *Point) Clone() *Point {
	if p == nil {
		return nil
	}
	return &Point{
		X:     p.X,
		Y:     p.Y,
		PType: p.PType,
	}
}
