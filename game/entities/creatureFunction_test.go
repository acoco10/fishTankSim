package entities

import (
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
	"github.com/acoco10/fishTankWebGame/game/geometry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"log"
	"testing"
)

func TestNormallyDistributedYRandomPoint(t *testing.T) {
	tankSize := geometry.Rect{X1: 0, Y1: 0, X2: 200, Y2: 200}
	fishStats, err := GenGoldFishStats()
	if err != nil {
		log.Fatal(err)
	}
	spriteSheet := spritesheet.SpriteSheet{SpriteWidth: 10}
	sp := sprite.AnimatedSprite{SpriteSheet: &spriteSheet}

	fish := Creature{FishStats: fishStats, TankBoundaries: tankSize, AnimatedSprite: &sp}

	nTests := 10000
	results := [10000]float32{}
	var sumResults float32
	for i := 0; i < nTests; i++ {
		targ := fish.RandomTarget()
		results[i] = targ.Y
		sumResults += targ.Y
	}

	t.Logf("avg Y: %f", sumResults/float32(10000))

}
