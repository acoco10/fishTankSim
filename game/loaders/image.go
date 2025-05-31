package loader

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func LoadImageAssetAsEbitenImage(assetName string) (*ebiten.Image, error) {
	imgPath := fmt.Sprintf("images/%s.png", assetName)
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, imgPath)
	if err != nil {
		return &ebiten.Image{}, err
	}
	return img, nil
}
