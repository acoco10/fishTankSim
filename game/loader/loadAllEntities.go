package loader

import (
	"fmt"
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

	z := rand.Intn(11) + 1
	zLayer := fmt.Sprintf("z%d", z)

	fish := entities.NewFishData(environment, hub, collisions[zLayer], fishName)
	newFishEntity := &entities.Entity{CreatureData: fish}
	println("loaded fish entity:", fishName.FishType)

	//mutates entity to have sprite and animations loaded
	entities.LoadFishSprite(newFishEntity)

	fp := getFirstPoint(collisions[zLayer])

	newFishEntity.Sprite.X = float32(fp.X)
	newFishEntity.Sprite.Y = float32(fp.Y)

	firstPoint := newFishEntity.RandomTarget(newFishEntity.CreatureData.MovementFlags)
	newFishEntity.MakeTargetPoint(firstPoint)
	newFishEntity.EventHub = hub

	entities.CreatureEventSubscriptions(newFishEntity)
	entities.RegisterEntity(newFishEntity)

	if registry.Config.Debug {
		newFishEntity.CreatureData.MaxHunger = 20
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

	return pt
}
