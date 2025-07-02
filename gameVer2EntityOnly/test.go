package main

import (
	entities2 "github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"log"
)

type Entity struct {
	Sprite *sprite.AnimatedSprite
	Active bool
}

func (e *Entity) Update() {

}

func Draw(Entiies []Entity, screen *ebiten.Image) {
	for _, e := range Entiies {
		if e.Sprite != nil {
			e.Sprite.Draw(screen)
		}
	}
}

type Game struct {
	Entities []Entity
}

func NewGame() *Game {
	sprite := loader.LoadFishSprite(entities2.Fish, 1)
	entities := make([]Entity, 100)
	for i := range entities {
		entities[i] = Entity{Sprite: sprite}
	}

	return &Game{
		Entities: entities,
	}
}

func (g *Game) Update() error {
	for _, i := range g.Entities {
		if i.Active {
			i.Sprite.Update()
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Black)
	for i := range g.Entities {
		if g.Entities[i].Active {
			g.Entities[i].Sprite.Draw(screen)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1280, 960
}

func main() {
	game := NewGame()

	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Generated Ebiten Skeleton")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
