package sprite

type Hovered struct {
	sprite *Sprite
}

func (h Hovered) Type() string {
	return "Hovered"
}
