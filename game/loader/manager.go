package loader

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/tasks"
)

type Manager struct {
	Hub *tasks.EventHub
}

func (m *Manager) Subscriptions() {
	m.Hub.Subscribe(entities.FishLevelUp{}, func(e tasks.Event) {
		ev := e.(entities.FishLevelUp)
		println("fish level up received for", ev.Fish.FishType)
		LoadLevlUpSprite(ev.Fish)
	})

}
