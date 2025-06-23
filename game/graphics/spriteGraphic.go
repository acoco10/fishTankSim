package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteGraphicEffect uint8

const (
	FadeIn SpriteGraphicEffect = iota
	Wipe
)

type SpriteGraphic struct {
	Sprite     sprite.Sprite
	updateFunc func(gs *SpriteGraphic)
	parameters map[string]any
	complete   bool
	drawFunc   func(gs *SpriteGraphic, screen *ebiten.Image)
}

func NewSpriteGraphic(gSprite sprite.Sprite, updateFunc func(gs *SpriteGraphic), params map[string]any) *SpriteGraphic {
	gs := SpriteGraphic{Sprite: gSprite, updateFunc: updateFunc, parameters: params}
	return &gs
}

func (gs *SpriteGraphic) Update() {
	gs.updateFunc(gs)
}

func (gs *SpriteGraphic) SetDrawFunc(effect SpriteGraphicEffect) {
	if effect == FadeIn {
		gs.drawFunc = DrawFadeInSprite
	}
}

func (gs *SpriteGraphic) Draw(screen *ebiten.Image) {
	if gs.drawFunc == nil {
		gs.Sprite.Draw(screen)
	} else {
		gs.drawFunc(gs, screen)
	}
}
