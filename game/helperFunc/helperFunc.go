package helperFunc

import (
	"fishTankWebGame/assets"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

func LoadImageAssetAsEbitenImage(assetName string) *ebiten.Image {
	imgPath := fmt.Sprintf("images/%s.png", assetName)
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, imgPath)
	if err != nil {
		log.Fatal(err)
	}
	return img
}
