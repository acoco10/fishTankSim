package main

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/scenes"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {

	loader.LoadShaderRegistry()
	loader.LoadFontRegistry()

	state := entities.SaveGameState{}
	gameLog := sceneManagement.NewGameLog(state, "")
	gameLog.Day = 1
	game := scenes.NewTestScene(sceneManagement.FishTank, gameLog)

	err := ebiten.RunGame(game)
	if err != nil {
		log.Fatal(err)
	}
}
