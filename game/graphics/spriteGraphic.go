package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteGraphicEffect uint8

const (
	FadeIn SpriteGraphicEffect = iota
	Wipe
)

type SpriteGraphic struct {
	Id         int
	Sprite     sprite.Sprite
	UpdateFunc func(gs *SpriteGraphic)
	Parameters map[string]any
	complete   bool
	drawFunc   func(gs *SpriteGraphic, screen *ebiten.Image)
}

type TopLevelSpriteGraphic struct {
	*SpriteGraphic
}

func (tls *TopLevelSpriteGraphic) Scaled() ScaledType {
	return ScaledButTopLevel
}

func (tls *TopLevelSpriteGraphic) Draw(screen *ebiten.Image) {
	tls.Sprite.Draw(screen)
}

func NewSpriteGraphic(gSprite sprite.Sprite, updateFunc func(gs *SpriteGraphic), params map[string]any) *SpriteGraphic {
	gs := SpriteGraphic{Sprite: gSprite, UpdateFunc: updateFunc, Parameters: params}
	gs.Id = AddGraphic(&gs)
	return &gs
}

func NewTopLevelSpriteGraphic(gSprite sprite.Sprite, updateFunc func(gs *SpriteGraphic), params map[string]any) *TopLevelSpriteGraphic {

	gSprite.Scale = registry.Config.ResolutionScalingF
	gs := &SpriteGraphic{Sprite: gSprite, UpdateFunc: updateFunc, Parameters: params}
	tgs := TopLevelSpriteGraphic{SpriteGraphic: gs}
	tgs.Id = AddGraphic(&tgs)
	return &tgs
}

func (gs *SpriteGraphic) AutoDeInit() bool {
	return true
}

func (gs *SpriteGraphic) Scaled() ScaledType {
	return NormalScaled
}

func (gs *SpriteGraphic) Update() {
	if gs.UpdateFunc == nil {
		return
	}
	gs.UpdateFunc(gs)
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
