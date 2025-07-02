package loader

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/props"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"image"
	"log"
)

func LoadProps(propList []entities.TankObject, tankBoundaries image.Rectangle, hub *tasks.EventHub) map[string]*props.StructureProp {
	propMap := make(map[string]*props.StructureProp)

	for _, prop := range propList {
		switch prop.Name {
		case "Log":
			logPropImg, err := LoadImageAssetAsEbitenImage("tankProps/logProp")
			logNormal, err := LoadImageAssetAsEbitenImage("tankProps/logProp_n")
			logProp := props.NewStructureProp(float32(tankBoundaries.Min.X), float32(tankBoundaries.Max.Y), logPropImg, logNormal, hub)
			propMap["Log"] = logProp
			if err != nil {
				log.Fatal(err)
			}
			logProp.StaticShadow = true
		case "Castle":
			castleImg, err := LoadImageAssetAsEbitenImage("tankProps/castleProp")
			castleNormal, err := LoadImageAssetAsEbitenImage("tankProps/castleProp_n")
			if err != nil {
				log.Fatal(err)
			}
			castleProp := props.NewStructureProp(float32(tankBoundaries.Min.X), float32(tankBoundaries.Max.Y), castleImg, castleNormal, hub)
			propMap["Castle"] = castleProp
		}
	}
	return propMap
}
