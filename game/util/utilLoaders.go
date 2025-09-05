package util

import (
	"bytes"
	"fmt"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
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
	case "childhood":
		loadedFont, err := assets.FontsDir.ReadFile("fonts/Childhood.otf")
		if err != nil {
			return nil, err
		}
		font = loadedFont
	case "miniPixel":
		loadedFont, err := assets.FontsDir.ReadFile("fonts/monogram.ttf")
		if err != nil {
			return nil, err
		}
		font = loadedFont
	case "ganon":
		loadedFont, err := assets.FontsDir.ReadFile("fonts/mini_pixel-7.ttf")
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
