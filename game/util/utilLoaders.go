package util

import (
	"bytes"
	"fmt"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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

func LoadImageAssetAsEbitenImage(assetName string) (*ebiten.Image, error) {
	imgPath := fmt.Sprintf("images/%s.png", assetName)
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, imgPath)
	if err != nil {
		return &ebiten.Image{}, err
	}
	return img, nil
}
