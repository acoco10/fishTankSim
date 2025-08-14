//go:build old

package scenes

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/registry"
)

var creatureManager CreatureManager

type CreatureManager struct {
	allFishFed bool
}


func SelectCreature(creature *entities.CreatureData) {
	creature.Selected = true
	creature.Shader = registry.ShaderMap["Outline"]
	loader.LoadRotatingHighlightOutlineAnimated(creature.AnimatedSprite)
}

func CheckIfAllFishFed(ents []*entities.Entity) bool {
	fed := true

	for _, ent := range ents {
		if ent.CreatureData != nil {
			if ent.CreatureData.Hunger > 0 {
				fed = false
			}
		}
	}
	return fed
}

func (g *FishScene2) creatureSubs(colMap map[string]image.Rectangle) {


	g.gameLog.GlobalEventHub.Subscribe(events.NewPurchase{}, func(e tasks.Event) {
		ev := e.(events.NewPurchase)
		log.Printf("New Purchase:%s ", ev.Purchase)
		creature := LoadPurchasedSprite(g.environment, ev.Purchase, g.gameLog.GlobalEventHub, g.collisionMap["tank"])
		newCreatureEntity := &entities.Entity{CreatureData: creature}
		entities.RegisterEntity(newCreatureEntity)
		g.gameEntities = append(g.gameEntities, newCreatureEntity)
	})

	g.gameLog.GlobalEventHub.Subscribe(entities.CreatureReachedPoint{}, func(e tasks.Event) {
		if !creatureManager.allFishFed && CheckIfAllFishFed(g.gameEntities) {
			creatureManager.allFishFed = true
			g.gameLog.GlobalEventHub.Publish(entities.AllFishFed{})
		}
	})
}
*/
