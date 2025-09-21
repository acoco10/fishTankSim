package registry

import (
	"bytes"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

var FontMap map[string]text.Face

const (
	PressStart   = "pressStart"
	PressStart12 = "pressStart_12"
	PressStart16 = "pressStart_16"
	PressStart24 = "pressStart_24"
	PressStart32 = "pressStart_32"
)

func LoadFontRegistry() {

	FontMap = make(map[string]text.Face)

	face, err := LoadFont(24, PressStart)
	if err != nil {
		log.Fatal("error loading pressstart", err)
	}
	FontMap["RockSalt"] = face

	face, err = LoadFont(20, PressStart)
	if err != nil {
		log.Fatal("error loading pressstart", err)
	}

	FontMap["RockSalt_18"] = face

	face, err = LoadFont(32, PressStart)
	if err != nil {
		log.Fatal("error loading pressstart", err)
	}
	FontMap["RockSalt_32"] = face

	face, err = LoadFont(12, PressStart)
	if err != nil {
		log.Fatal("error loading pressstart", err)
	}
	FontMap["RockSalt_12"] = face

	face, err = LoadFont(16, PressStart)
	if err != nil {
		log.Fatal("error loading pressstart", err)
	}
	FontMap["RockSalt_16"] = face

	face, err = LoadFont(16, "nk57")
	if err != nil {
		log.Fatal("error loading pressstart", err)
	}
	FontMap["nk57"] = face

	face, err = LoadFont(12, "nk57")
	if err != nil {
		log.Fatal("error loading pressstart", err)
	}
	FontMap["nk57_12"] = face

	face, err = LoadFont(24, "nk57")
	if err != nil {
		log.Fatal("error loading nk56", err)
	}
	FontMap["nk57_24"] = face

}

const (
	Nk57      = "nk57"
	RockSalt  = "rockSalt"
	StayPixel = "stayPixel"
)

func LoadFont(size float64, fontName string) (text.Face, error) {
	var font []byte
	switch fontName {
	case "nk57":
		loadedFont, err := assets.FontsDir.ReadFile("fonts/nk57.otf")
		if err != nil {
			return nil, err
		}
		font = loadedFont
	case "pressStart":
		loadedFont, err := assets.FontsDir.ReadFile("fonts/PressStart2P.ttf")
		if err != nil {
			return nil, err
		}
		font = loadedFont
	case "stayPixel":
		loadedFont, err := assets.FontsDir.ReadFile("fonts/stayPixel.ttf")
		if err != nil {
			return nil, err
		}
		font = loadedFont
	case "reglisseOutlined":
		loadedFont, err := assets.FontsDir.ReadFile("fonts/Reglisse_Outlined.otf")
		if err != nil {
			return nil, err
		}
		font = loadedFont
	case "reglisseOutline":
		loadedFont, err := assets.FontsDir.ReadFile("fonts/Reglisse_Outline.otf")
		if err != nil {
			return nil, err
		}
		font = loadedFont
	case "reglisseClean":
		loadedFont, err := assets.FontsDir.ReadFile("fonts/Reglisse_Clean.otf")
		if err != nil {
			return nil, err
		}
		font = loadedFont
	}
	s, err := text.NewGoTextFaceSource(bytes.NewReader(font))
	if err != nil {
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}
