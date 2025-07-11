package interactableUIObjects

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type EventUI struct {
	*UiSprite
	Triggered bool
}

func (p *EventUI) Update() {

	if p.SpriteHovered() {
		//shader := shaders.LoadOutlineShader()
		//p.Shader = shader
		//p.ShaderParams["OutlineColor"] = []float64{1, 1, 0, 1}
	}

	if p.Shader != nil && !p.SpriteHovered() {
		//p.Shader = nil
	}

	if p.Triggered && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && p.SpriteHovered() {
		p.EventHub.Publish(events.UISpriteAction{UiSprite: p.Label, UiSpriteAction: "clicked"})
		turnOffClickMeEffect(p.UiSprite)
	}

	p.Sprite.Update()

}

func (p *EventUI) Draw(screen *ebiten.Image) {
	if p.Triggered == true {
		p.Sprite.Draw(screen)
	}
}

func (p *EventUI) Subscribe(hub *tasks.EventHub) {
	if p.Label == "pillow" {
		hub.Subscribe(tasks.AllTasksCompleted{}, func(e tasks.Event) {
			initClickMeEffect(p.UiSprite)
			p.Triggered = true
		})
	}

	if p.Label == "door" {
		hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
			ev := e.(events.NewDay)
			if ev.Type == "Chores" {
				initClickMeEffect(p.UiSprite)
				p.Triggered = true
			}
		})
	}

	hub.Subscribe(events.DayOver{}, func(e tasks.Event) {
		p.Triggered = false
	})
}
