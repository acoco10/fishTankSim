package util

import (
	"bytes"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
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
	case "rockSalt":
		loadedFont, err := assets.FontsDir.ReadFile("fonts/RockSalt.ttf")
		if err != nil {
			return nil, err
		}
		font = loadedFont
	}
	s, err := text.NewGoTextFaceSource(bytes.NewReader(font))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}
