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

func LoadAllEntities(uiSpritesToLoad []entities.Label, fishList []entities.SavedFish, environment *system.Environment, hub *tasks.EventHub, collisions map[string]image.Rectangle, gs entities.GameState) ([]*entities.Entity, *entities.WhiteBoardSprite) {

	var entitySlice []*entities.Entity
	var wbSprite *entities.WhiteBoardSprite

	for _, fishName := range fishList {
		ent := InitFish(fishName, environment, hub, gs)
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

func InitFish(fishName entities.SavedFish, environment *system.Environment, hub *tasks.EventHub, state entities.GameState) *entities.Entity {

	z := rand.Intn(10) + 2

	fish := entities.NewFishData(environment, hub, state.Zbounds[z], fishName)
	fish.TargetZ = z
	newFishEntity := &entities.Entity{CreatureData: fish, Z: z}
	println("loaded fish entity:", fishName.FishType)

	//mutates entity to have sprite and animations loaded
	entities.LoadFishSprite(newFishEntity)
	newFishEntity.Sprite.DOptsUpdaterTag = "custom"

	fp := getFirstPoint(state.Zbounds[z])

	newFishEntity.Sprite.X = float32(fp.X)
	newFishEntity.Sprite.Y = float32(fp.Y)

	firstPoint := newFishEntity.RandomTarget()
	newFishEntity.SetTargetPoint(firstPoint)
	newFishEntity.EventHub = hub
	newFishEntity.FishUpdate(state)
	newFishEntity.Sprite.Update() //update the animations draw opts so it doesnt draw at 0,0 lolololol

	entities.CreatureEventSubscriptions(newFishEntity)
	entities.RegisterEntity(newFishEntity)

	if registry.Config.Debug {
		newFishEntity.CreatureData.MaxHunger = 20
	}

	return newFishEntity
}

func getFirstPoint(bounds image.Rectangle) image.Point {
	forgivingRect := bounds
	forgivingRect.Min.X += 50
	forgivingRect.Min.Y += 50
	forgivingRect.Max.X -= 50
	forgivingRect.Max.Y -= 50

	firstPointX := forgivingRect.Min.X + rand.Intn(forgivingRect.Max.X-forgivingRect.Min.X)
	firstPointY := forgivingRect.Min.Y + rand.Intn(forgivingRect.Max.Y-forgivingRect.Min.Y)
	pt := image.Point{firstPointX, firstPointY}

	return pt
}
