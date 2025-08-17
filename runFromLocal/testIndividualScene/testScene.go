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

	//loader.LoadShaderRegistry()
	loader.LoadFontRegistry()

	state := entities.SaveGameState{}
	gameLog := sceneManagement.NewGameLog(state, "w")

	game := scenes.NewTestScene(sceneManagement.MowingMiniGameScene, gameLog)

	err := ebiten.RunGame(game)
	if err != nil {
		log.Fatal(err)
	}
}
