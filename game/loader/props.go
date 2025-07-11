package loader

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/props"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"image"
	"log"
)

func LoadProps(propList []entities.TankObject, tankBoundaries image.Rectangle, hub *tasks.EventHub) *props.PropQueue {
	pQueue := props.PropQueue{}

	for _, prop := range propList {
		switch prop.Name {
		case "Log":
			logPropImg, err := LoadImageAssetAsEbitenImage("tankProps/logProp")
			logNormal, err := LoadImageAssetAsEbitenImage("tankProps/logProp_n")
			logProp := props.NewStructureProp(float32(tankBoundaries.Min.X), float32(tankBoundaries.Max.Y), logPropImg, logNormal, hub)

			if err != nil {
				log.Fatal(err)
			}
			pQueue.ActiveProp = logProp
		case "Castle":
			castleImg, err := LoadImageAssetAsEbitenImage("tankProps/castleProp")
			castleNormal, err := LoadImageAssetAsEbitenImage("tankProps/castleProp_n")
			if err != nil {
				log.Fatal(err)
			}
			castleProp := props.NewStructureProp(float32(tankBoundaries.Min.X), float32(tankBoundaries.Max.Y), castleImg, castleNormal, hub)
			pQueue.ActiveProp = castleProp
		}
	}

	grassImg, err := LoadImageAssetAsEbitenImage("tankProps/grass")
	if err != nil {
		log.Fatal(err, "tried to load plant img from wrong place")
	}
	grassNormal, err := LoadImageAssetAsEbitenImage("tankProps/grass_n")
	if err != nil {
		log.Fatal(err)
	}
	grassProp := props.NewStructureProp(float32(tankBoundaries.Min.X), float32(tankBoundaries.Max.Y), grassImg, grassNormal, hub)
	pQueue.QueuedProps = append(pQueue.QueuedProps, grassProp)

	return &pQueue
}
