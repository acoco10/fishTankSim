package sprite

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/util"
	"image"
	"image/color"
)

func (s *Sprite) AddColoredOutlineShader(rgba color.RGBA) {

	s.ShaderParams["OutlineColor"] = util.ConvertRGBAUintToDecimal(rgba)
	if s.NormalMap != nil {
		s.Shader = registry.ShaderMap[registry.NormalMapOutline]
	}
	s.Shader = registry.ShaderMap[registry.Outline]
}

func (s *Sprite) ChangeOutlineColor(rgba color.RGBA) {
	s.ShaderParams["OutlineColor"] = util.ConvertRGBAUintToDecimal(rgba)
}

func (s *Sprite) MidX() int {
	return s.GetSpriteRect().Dx() / 2
}

func (s *Sprite) TranslatedMidX() int {
	return int(s.X) + s.GetSpriteRect().Dx()/2
}

func (s *Sprite) MidPoint() (int, int) {
	return s.GetSpriteRect().Dx() / 2, s.GetSpriteRect().Dy() / 2
}

func (s *Sprite) TranslatedMidPointAsPoint() image.Point {
	return image.Point{s.TranslatedMidX(), s.TranslatedMidY()}
}

func (s *Sprite) MidY() int {
	return s.GetSpriteRect().Dy() / 2
}

func (s *Sprite) TranslatedMidY() int {
	return int(s.Y) + s.GetSpriteRect().Dy()/2
}

func (s *Sprite) CheckMiddleOfSpriteInRect(rect image.Rectangle) bool {
	return s.TranslatedMidPointAsPoint().In(rect)
}
