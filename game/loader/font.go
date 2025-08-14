package loader

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

func LoadFontRegistry() {

	registry.FontMap = make(map[string]text.Face)

	face, err := util.LoadFont(16, "rockSalt")
	if err != nil {
		log.Fatal(err)
	}
	registry.FontMap["RockSalt"] = face

	face, err = util.LoadFont(12, "rockSalt")
	if err != nil {
		log.Fatal(err)
	}
	registry.FontMap["RockSalt_12"] = face

	face, err = util.LoadFont(16, "nk57")
	if err != nil {
		log.Fatal(err)
	}
	registry.FontMap["nk57"] = face

	face, err = util.LoadFont(12, "nk57")
	if err != nil {
		log.Fatal(err)
	}
	registry.FontMap["nk57_12"] = face

	face, err = util.LoadFont(24, "nk57")
	if err != nil {
		log.Fatal(err)
	}
	registry.FontMap["nk57_24"] = face

}
