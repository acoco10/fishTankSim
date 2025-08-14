package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

func NewVlS(x0, y0, x1, y1, maxX float32, clr color.Color, dst *ebiten.Image) int {
	vls := &VectorLineGraphic{
		X0:       x0,
		Y0:       y0,
		X1:       x1,
		Y1:       y1,
		Color:    clr,
		MaxX1:    maxX,
		dstImage: dst,
	}

	id := AssignAndIncrement(vls)
	return id
}

func LoadClothGraphic(origin [2]float32) (*TopLevelSpriteGraphic, error) {

	origin[0] *= registry.Config.ResolutionScalingf
	origin[1] *= registry.Config.ResolutionScalingf
	origin[1] += float32(registry.Config.ScaledYOffsetF)
	mx := [2]float32{origin[0] + 150, origin[1] + 150}

	direction := "right"

	img, err := util.LoadImageAssetAsEbitenImage("menuAssets/cloth")
	if err != nil {
		return nil, err
	}
	gSprite := sprite.Sprite{Img: img, X: origin[0], Y: origin[1]}

	params := make(map[string]any)
	params["origin"] = origin
	params["max"] = mx
	params["direction"] = direction
	cloth := &SpriteGraphic{Sprite: gSprite, UpdateFunc: UpdateWhiteBoardCloth, Parameters: params}
	tls := &TopLevelSpriteGraphic{SpriteGraphic: cloth}
	tls.Sprite.Scale = registry.Config.ResolutionScalingF
	return tls, nil
}
