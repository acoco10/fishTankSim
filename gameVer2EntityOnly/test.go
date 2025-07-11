package main

import (
	"github.com/acoco10/fishTankWebGame/gameVer2EntityOnly/entity"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"log"
)

type Game struct {
	Entities []entity.Entity
	Plant    *entity.Plant
}

func NewGame() *Game {

	plant := entity.MakePlant(10)

	return &Game{
		Plant: plant,
	}
}

func (g *Game) Update() error {
	g.Plant.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Black)
	g.Plant.Draw(screen)
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
