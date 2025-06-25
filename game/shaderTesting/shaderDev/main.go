package main

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/interactableUIObjects"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	"log"
)

const (
	screenWidth  = 640 * 2
	screenHeight = 480 * 2
)

type Game struct {
	eventHub   *tasks.EventHub
	whiteBoard *interactableUIObjects.WhiteBoardSprite2
	testTask   *tasks.Task
}

func newGame() *Game {

	g := Game{}
	hub := tasks.NewEventHub()
	g.eventHub = hub

	taskCondition2 := func(e tasks.Event) bool {
		ev, ok := e.(entities.SendData)
		return ok && ev.DataFor == "statsMenu"
	}

	gameTask2 := tasks.NewTask(entities.SendData{}, "2. Click your fish", taskCondition2)

	gameTask2.Subscribe(g.eventHub)

	g.testTask = gameTask2

	//shader := shaders
	//s.shader = outlineShader

	wb := interactableUIObjects.WhiteBoardSprite2{}
	wbImg, err := loader.LoadImageAssetAsEbitenImage("uiSprites/whiteBoardMain")
	if err != nil {
		log.Fatal(err)
	}

	wbuiSprte := &interactableUIObjects.UiSprite{}
	wb.UiSprite = wbuiSprte

	wbSprite := &sprite.Sprite{Img: wbImg, X: 300, Y: 300}
	wb.Sprite = wbSprite
	wb.EventHub = g.eventHub
	wb.Init()
	g.whiteBoard = &wb

	return &g
}

func (g *Game) Update() error {
	g.whiteBoard.Update()

	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.testTask.Activate()
		g.testTask.Publish(g.eventHub)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(colornames.Darkgreen)
	g.whiteBoard.Draw(screen)

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
