package main

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

const (
	screenWidth  = 1000
	screenHeight = 1000
)

type Game struct {
	testGraphic *graphics.SpriteGraphic
}

func newGame() *Game {
	g := Game{}

	graphSpriteParams := map[string]any{}

	graphSpriteParams["opacity"] = float32(0.1)
	graphSpriteParams["pulse"] = true

	sprite := loader.LoadFishSprite(entities.Fish, 2)
	sprite.X = 200
	sprite.Y = 200
	graph := graphics.NewSpriteGraphic(*sprite.Sprite, graphics.UpdateFadeInGraphic, graphSpriteParams)
	graph.SetDrawFunc(graphics.FadeIn)
	g.testGraphic = graph
	return &g

}

func (g *Game) Update() error {
	g.testGraphic.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.testGraphic.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1000, 1000
}

func main() {
	g := newGame()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Graphic Testing")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal("error running game", err)
	}
}
