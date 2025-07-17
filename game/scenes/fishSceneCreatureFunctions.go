package scenes

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/geometry"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/input"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
	"math/rand"
)

type CreatureManager struct {
	allFishFed bool
}

func (g *FishScene) checkForFishSelected() {
	if g.currentTask > 0 || g.gameLog.Day > 1 {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !g.mouseFlags.CursorOccupied {
			x, y := ebiten.CursorPosition()
			xCheck := x > g.tankSize.Min.X && x < g.tankSize.Max.X
			yCheck := y > g.tankSize.Min.Y && y < g.tankSize.Max.Y

			if xCheck && yCheck {
				filterFunc := func(distance any) bool {
					return distance.(float64) < 50
				}
				cursorX, cursorY := ebiten.CursorPosition()
				closestCreature := util.ClosestDrawableToCursor(cursorX, cursorY, g.sprites[0], filterFunc, "*entities.Creature")

				if closestCreature != nil {
					cre, ok := closestCreature.(*entities.Creature)
					if ok {
						SelectCreature(cre)
						AddEffectToClosestCreature(g.lightingState, cre)
					}
				}
			}
		}
	}
}

func AddEffectToClosestCreature(state lightingState, cre *entities.Creature) {
	switch state {
	case Night:
		sleepyEffect := loader.LoadEffect("zzz")
		graphics.NewTravelingEffect(sleepyEffect, &cre.X, &cre.Y)
	case Day:
		Effect := loader.LoadEffect("happy")
		graphics.NewTravelingEffect(Effect, &cre.X, &cre.Y)
	}
}

func SelectCreature(creature *entities.Creature) {
	creature.Selected = true
	creature.Shader = registry.ShaderMap["Outline"]
	loader.LoadRotatingHighlightOutlineAnimated(creature.AnimatedSprite)
}

func (g *FishScene) CheckIfAllFishFed() bool {
	fed := true

	for _, draw := range g.sprites[0] {
		creature, ok := draw.(*entities.Creature)
		if ok && creature.Hunger != 0 {
			fed = false
		}
	}
	return fed

}

func (g *FishScene) creatureSubs(colMap map[string]geometry.Rect) {
	g.gameLog.GlobalEventHub.Subscribe(input.MouseButtonPressedUISpriteActivity{}, func(e tasks.Event) {
		ev := e.(input.MouseButtonPressedUISpriteActivity)
		if g.timers["pointGeneratedTimer"].TimerState == entities.Done && !g.mouseFlags.HandledClick {
			g.mouseFlags.HandledClick = true
			pt := ev.Point.Clone()
			pt.X = pt.X - 50 + rand.Float32()*10
			pt.Y += 50
			p := entities.NewParticle(pt, colMap["tank"], g.gameLog.GlobalEventHub)
			g.sprites[0] = append(g.sprites[0], p)
		}
	})

	g.gameLog.GlobalEventHub.Subscribe(events.NewPurchase{}, func(e tasks.Event) {
		ev := e.(events.NewPurchase)
		log.Printf("New Purchase:%s ", ev.Purchase)
		creature := LoadPurchasedSprite(g.environment, ev.Purchase, g.gameLog.GlobalEventHub, g.collisionMap["tank"])
		g.sprites[0] = append(g.sprites[0], creature)
	})

}
