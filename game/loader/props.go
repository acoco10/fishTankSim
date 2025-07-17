package loader

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/props"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"image"
	"log"
)

func LoadProp(propName string, tankBoundaries image.Rectangle, eventhub *tasks.EventHub) *props.StructureProp {
	switch propName {
	case "Log":
		logPropImg, err := LoadImageAssetAsEbitenImage("tankProps/logProp")
		logNormal, err := LoadImageAssetAsEbitenImage("tankProps/logProp_n")
		logProp := props.NewStructureProp(0, 0, logPropImg, logNormal, eventhub, tankBoundaries)

		if err != nil {
			log.Fatal(err)
		}
		return logProp
	case "Castle":

		log.Println("returning castle prop from load prop call")
		castleImg, err := LoadImageAssetAsEbitenImage("tankProps/castleProp")
		castleNormal, err := LoadImageAssetAsEbitenImage("tankProps/castleProp_n")
		if err != nil {
			log.Fatal(err)
		}
		castleProp := props.NewStructureProp(0, 0, castleImg, castleNormal, eventhub, tankBoundaries)
		return castleProp

	case "Grass":
		grassImg, err := LoadImageAssetAsEbitenImage("tankProps/grass")
		if err != nil {
			log.Fatal(err, "tried to load plant img from wrong place")
		}

		grassNormal, err := LoadImageAssetAsEbitenImage("tankProps/grass_n")
		if err != nil {
			log.Fatal(err)
		}

		grassProp := props.NewStructureProp(0, 0, grassImg, grassNormal, eventhub, tankBoundaries)
		return grassProp
	}
	return nil
}

func LoadProps(propList []entities.TankObject, tankBoundaries image.Rectangle, hub *tasks.EventHub) props.PropQueue {
	pQueue := props.PropQueue{}

	for _, prop := range propList {
		p := LoadProp(prop.Name, tankBoundaries, hub)
		pQueue.QueuedProps = append(pQueue.QueuedProps, p)
	}

	return pQueue
}
