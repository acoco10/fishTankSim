package loader

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/system"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"image"
	"log"
	"math/rand"
)

func LoadAllEntities(uiSpritesToLoad []entities.Label, fishList []entities.SavedFish, environment *system.Environment, hub *tasks.EventHub, collisions map[string]image.Rectangle) ([]*entities.Entity, *entities.WhiteBoardSprite) {

	var entitySlice []*entities.Entity
	var wbSprite *entities.WhiteBoardSprite

	for _, fishName := range fishList {
		ent := InitFish(fishName, environment, hub, collisions)
		entitySlice = append(entitySlice, ent)
	}

	uiSprites, wb, err := LoadUISprites(uiSpritesToLoad, environment, collisions["tank"], hub)
	if err != nil {
		log.Fatal("Error loading uisprites from load entities function", err)
	}

	wbSprite = wb
	entitySlice = append(entitySlice, uiSprites...)

	return entitySlice, wbSprite
}

func InitFish(fishName entities.SavedFish, environment *system.Environment, hub *tasks.EventHub, collisions map[string]image.Rectangle) *entities.Entity {
	fish := entities.NewFishData(environment, hub, collisions["tank"], fishName)
	sprite, err := entities.LoadFishSprite(fish.FishType, fish.Size)

	sprite.X = float32(collisions["tank"].Min.X) + rand.Float32()*float32(collisions["tank"].Max.X)
	sprite.Y = float32(collisions["tank"].Min.Y) + rand.Float32()*float32(collisions["tank"].Max.Y)
	sprite.Z = 1
	sprite.Unfocusable = true

	if err != nil {
		log.Fatal(err)
	}
	newFishEntity := &entities.Entity{CreatureData: fish, Sprite: sprite}
	println("loaded fish entity:", fishName.FishType)

	firstPoint := newFishEntity.RandomTarget()
	newFishEntity.MakeTargetPoint(firstPoint)
	newFishEntity.EventHub = hub

	entities.CreatureEventSubscriptions(newFishEntity)

	entities.RegisterEntity(newFishEntity)

	return newFishEntity
}
