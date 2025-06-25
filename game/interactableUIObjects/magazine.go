package interactableUIObjects

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//import "github.com/acoco10/fishTankWebGame/game/events"

type MagazineUI struct {
	*UiSprite
	//ui *ui.Magazine
}

func (m *MagazineUI) Update() {
	if m.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ev := events.UISpriteAction{
			UiSprite:       "magazine",
			UiSpriteAction: "picked up",
		}
		m.EventHub.Publish(ev)
	}
}
