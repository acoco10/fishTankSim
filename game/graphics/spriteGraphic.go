package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteGraphic struct {
	Sprite     *sprite.Sprite
	updateFunc func(gs *SpriteGraphic)
	parameters map[string]any
	complete   bool
}

func NewSpriteGraphic(gSprite *sprite.Sprite, updateFunc func(gs *SpriteGraphic), params map[string]any) *SpriteGraphic {
	gs := SpriteGraphic{Sprite: gSprite, updateFunc: updateFunc, parameters: params}
	return &gs
}

func (gs *SpriteGraphic) Update() {
	gs.updateFunc(gs)
}

func (gs *SpriteGraphic) Draw(screen *ebiten.Image) {
	gs.Sprite.Draw(screen)
}
