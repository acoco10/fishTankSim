package util

import (
	"bytes"
	"fmt"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"
)

func LoadImageAssetAsEbitenImage(assetName string) (*ebiten.Image, error) {
	imgPath := fmt.Sprintf("images/%s.png", assetName)
	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, imgPath)
	if err != nil {
		return &ebiten.Image{}, err
	}
	return img, nil
}

func LoadDirectoryImages(directory string) map[string]*ebiten.Image {

	dir, err := assets.ImagesDir.ReadDir(directory)

	if err != nil {
		log.Printf("Error reading directory: %v", err)
		return nil
	}

	imageMap := make(map[string]*ebiten.Image)

	for _, file := range dir {
		if file.IsDir() {
			continue // Skip subdirectories
		}

		// Read the file
		data, err := assets.ImagesDir.ReadFile(directory + "/" + file.Name())
		if err != nil {
			log.Printf("Error reading file %s: %v", file.Name(), err)
			continue // Skip files that can't be read
		}

		// Decode the image
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			log.Printf("Error decoding image %s: %v", file.Name(), err)
			continue // Skip files that can't be decoded as images
		}

		// Convert to ebiten image
		ebitenImg := ebiten.NewImageFromImage(img)

		// Use filename as key (you might want to remove extension)
		fileName := file.Name()[:len(file.Name())-4]
		log.Printf("Successfully loaded image: %s from directory %s", fileName, directory)

		imageMap[fileName] = ebitenImg
	}

	return imageMap
}
