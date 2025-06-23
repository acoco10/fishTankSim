package ui

import (
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

type Manager struct {
	uiMap map[string]*ebitenui.UI
	uiKey string
}

func LoadFishSceneUIManager(hub *tasks.EventHub, width int, height int) *Manager {
	m := Manager{}

	mainUI, _, err := LoadMainFishMenu(width, height, hub)
	if err != nil {
		log.Fatal(err)
	}

	uiMap := make(map[string]*ebitenui.UI)

	uiMap["main"] = mainUI

	return &m
}

func (m *Manager) Update() {
	m.uiMap[m.uiKey].Update()
}

func (m *Manager) Draw(screen *ebiten.Image) {
	m.uiMap[m.uiKey].Draw(screen)
}
