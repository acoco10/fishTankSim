package entities

import (
	"github.com/acoco10/fishTankWebGame/game/entImportableLoaders"
	"github.com/acoco10/fishTankWebGame/game/sprite"
)

func LoadStomachGraphic() *sprite.Sprite {
	sp := entImportableLoaders.LoadStaticEffect("stomach", 0, 0, "")
	return sp
}
