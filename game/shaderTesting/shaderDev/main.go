package main

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/input"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/system"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

const (
	screenWidth  = 640 * 2
	screenHeight = 480 * 2
)

type Game struct {
	eventHub    *tasks.EventHub
	gameState   *entities.GameState
	environment *system.Environment
	fishEnt     *entities.Entity
}

func newGame() *Game {

	g := Game{}
	hub := tasks.NewEventHub()
	g.eventHub = hub
	mouseFlags := &input.MouseFlags{HandledClick: false, CursorOccupied: false}
	g.gameState = &entities.GameState{MouseFlags: mouseFlags}
	g.environment = &system.Environment{}
	system.InitEnvironment(g.environment)
	registry.Config.Zoom = false
	registry.Config.ZoomFactor = 2
	registry.Config.ResolutionHeight = 800
	registry.Config.ResolutionWidth = 480 * 2

	fish := entities.SavedFish{FishType: string(entities.Fish), Size: 1}

	colMap := make(map[string]image.Rectangle)
	colMap["tank"] = image.Rect(100, 100, 400, 400)

	fEnt := loader.InitFish(fish, g.environment, g.eventHub, colMap)
	g.fishEnt = fEnt

	return &g
}

func (g *Game) Update() error {
	if g.fishEnt.Sprite.Focused == false {
		entities.Focus(g.fishEnt.Id)
	}
	entities.UpdateEntities(g.gameState)
	return nil

}

func (g *Game) Draw(screen *ebiten.Image) {

	entities.DrawEntities(screen)
	entities.DrawNonZoomedEntities(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	loader.LoadFontRegistry()
	loader.LoadShaderRegistry()
	g := newGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hand writing shader")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
