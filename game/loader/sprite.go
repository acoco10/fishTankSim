package loader

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
)

func MakeSpriteEntity(img *ebiten.Image, x float32, y float32, flags SpriteEntFlags) uint32 {

	sp := &sprite.Sprite{Img: img, X: x, Y: y}
	sp.Unfocusable = flags.Unfocusable
	ent := &entities.Entity{Sprite: sp}
	entities.RegisterEntity(ent)

	if flags.Updater {
		sp.XYUpdater = sprite.NewUpdater(sp)
	}

	ent.UpdateFunc = flags.UpdateFunc
	ent.Parameters = flags.Parameters
	ent.EventHub = flags.EventHub
	ent.Z = flags.Zlayer

	return ent.Id
}

type SpriteEntFlags struct {
	Unfocusable bool
	Zlayer      int
	Updater     bool
	UpdateFunc  func(ent *entities.Entity)
	Parameters  map[string]any
	EventHub    *tasks.EventHub
	EventType   tasks.Event
	sub         func(event tasks.Event)
}
