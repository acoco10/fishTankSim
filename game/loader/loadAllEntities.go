package loader

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/registry"
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
	aniMap, err := entities.LoadFishSprite(fish.FishType, fish.Size)

	sprite := aniMap["swimming"]

	sprite.Z = 1
	sprite.Unfocusable = true

	fp := getFirstPoint(collisions["tank"])

	sprite.X = float32(fp.X)
	sprite.Y = float32(fp.Y)

	if err != nil {
		log.Fatal(err)
	}
	newFishEntity := &entities.Entity{CreatureData: fish, Sprite: sprite}
	println("loaded fish entity:", fishName.FishType)

	firstPoint := newFishEntity.RandomTarget()
	newFishEntity.MakeTargetPoint(firstPoint)
	newFishEntity.EventHub = hub
	newFishEntity.AnimationMap = aniMap
	entities.CreatureEventSubscriptions(newFishEntity)

	entities.RegisterEntity(newFishEntity)

	if registry.Config.Debug {
		//newFishEntity.CreatureData.MaxHunger = 20
	}

	return newFishEntity
}

func getFirstPoint(bounds image.Rectangle) image.Point {
	forgivingRect := bounds
	forgivingRect.Min.X += 5
	forgivingRect.Min.Y += 5
	forgivingRect.Max.X -= 5
	forgivingRect.Min.X += 5

	firstPointX := forgivingRect.Min.X + 5 + rand.Intn(bounds.Max.X-5-bounds.Min.X+5)
	firstPointY := forgivingRect.Min.Y + 5 + rand.Intn(bounds.Max.X-5-bounds.Min.Y+5)
	pt := image.Point{firstPointX, firstPointY}

	if !pt.In(forgivingRect) {
		pt = getFirstPoint(bounds)
	}

	return pt
}
