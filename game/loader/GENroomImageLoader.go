package loader

import (
	"bytes"
	"fmt"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type BackGroundImages struct {
	OffScreen                        *ebiten.Image
	OffScreen2                       *ebiten.Image
	RoomBackground                   *ebiten.Image
	FishTank                         *ebiten.Image
	FishTankDayLight                 *ebiten.Image
	FishTank_n                       *ebiten.Image
	FishTankFrontLayerDayLight       *ebiten.Image
	FishTankNight                    *ebiten.Image
	FishTankFrontLayerNoLightSmaller *ebiten.Image
	LaptopRoomBackground             *ebiten.Image
	RoomBackGroundNight              *ebiten.Image
	FishTankFrontLayerZoom           *ebiten.Image
}

func LoadImageFromEmbedded(directory, filename string) (*ebiten.Image, error) {
	img, err := util.LoadImageAssetAsEbitenImage(directory + filename)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func LoadAllRoomBackGroundImages(directory string) (*BackGroundImages, error) {
	var m BackGroundImages
	var err error

	m.RoomBackground, err = LoadImageFromEmbedded(directory, "roomBackground")
	if err != nil {
		return nil, fmt.Errorf("loading roomBackground: %w", err)
	}

	m.FishTank, err = LoadImageFromEmbedded(directory, "fishTank")
	if err != nil {
		return nil, fmt.Errorf("loading fishTank: %w", err)
	}

	m.FishTankDayLight, err = LoadImageFromEmbedded(directory, "fishTankDayLight")
	if err != nil {
		return nil, fmt.Errorf("loading fishTankDayLight: %w", err)
	}

	m.FishTank_n, err = LoadImageFromEmbedded(directory, "fishTank_n")
	if err != nil {
		return nil, fmt.Errorf("loading fishTank_n: %w", err)
	}

	m.FishTankFrontLayerDayLight, err = LoadImageFromEmbedded(directory, "fishTankFrontLayerDayLight")
	if err != nil {
		return nil, fmt.Errorf("loading fishTankFrontLayerDayLight: %w", err)
	}

	m.FishTankNight, err = LoadImageFromEmbedded(directory, "fishTankNight")
	if err != nil {
		return nil, fmt.Errorf("loading fishTankNight: %w", err)
	}

	m.FishTankFrontLayerNoLightSmaller, err = LoadImageFromEmbedded(directory, "fishTankFrontLayerNoLightSmaller")
	if err != nil {
		return nil, fmt.Errorf("loading fishTankFrontLayerNoLightSmaller: %w", err)
	}

	m.LaptopRoomBackground, err = LoadImageFromEmbedded(directory, "laptopRoomBackground")
	if err != nil {
		return nil, fmt.Errorf("loading laptopRoomBackground: %w", err)
	}

	m.RoomBackGroundNight, err = LoadImageFromEmbedded(directory, "roomBackgroundNight")
	if err != nil {
		return nil, fmt.Errorf("loading roomBackGroundNight: %w", err)
	}

	m.FishTankFrontLayerZoom, err = LoadImageFromEmbedded(directory, "fishTankFrontLayerZoom")
	if err != nil {
		return nil, fmt.Errorf("loading fishTankFrontLayerZoom: %w", err)
	}

	return &m, nil
}

// Keep your existing utility functions unchanged
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
