package entities

type entFlags uint32

const (
	overUi = 1 << iota
)

func (ent *Entity) SetOverUIEffect() {
	ent.flags |= overUi
}

func (ent *Entity) IsOverUI() bool {
	return ent.flags&overUi != 0
}
